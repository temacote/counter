package cmd

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/philips/go-bindata-assetfs"
	"github.com/spf13/cobra"
	"go.uber.org/zap"

	"sber_cloud/tw/cmd/definition/http/gateway"
	middleware2 "sber_cloud/tw/cmd/definition/http/middleware"
	"sber_cloud/tw/cmd/http/middleware"
	counter2 "sber_cloud/tw/counter"
	counter_def "sber_cloud/tw/definition/counter"
	"sber_cloud/tw/definition/logger"
	"sber_cloud/tw/definition/monitoring"
	redis_def "sber_cloud/tw/definition/redis"
	"sber_cloud/tw/redis"
	"sber_cloud/tw/swagger"
)

var (
	httpCmd = &cobra.Command{
		Use:   "http",
		Short: "Run HTTP server",
		RunE:  httpServe,
	}
	swaggerPath string
)

// Command init function.
func init() {
	httpCmd.PersistentFlags().StringVarP(&swaggerPath, "swagger", "w", "swagger.json", "enable swagger")
	rootCmd.AddCommand(httpCmd)
}

func httpServe(_ *cobra.Command, _ []string) (err error) {
	var httpListen = conf.GetString("http.listen")

	var log logger.Logger
	if err = diContainer.Fill(logger.DefLogger, &log); err != nil {
		return err
	}

	var mux *http.ServeMux
	if err = diContainer.Fill(gateway.DefHTTPGateway, &mux); err != nil {
		return err
	}

	var swaggerEnabled bool
	func() {
		if len(swaggerPath) == 0 {
			return
		}
		if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
			log.Info("swagger file not found", zap.String("path", swaggerPath))
			return
		}

		var swaggerFile *os.File
		if swaggerFile, err = os.Open(swaggerPath); err != nil {
			log.Error("can't open swagger.json", zap.String("path", swaggerPath), zap.Error(err))
			return
		}

		var swaggerData []byte
		if swaggerData, err = ioutil.ReadAll(swaggerFile); err != nil {
			log.Error("error read swagger.json", zap.String("path", swaggerPath), zap.Error(err))
		}
		mux.HandleFunc("/swagger.json", func(w http.ResponseWriter, req *http.Request) {
			if _, err = w.Write(swaggerData); err != nil {
				log.Error("error send swagger.json data")
			}
		})
		serveSwagger(mux)
		swaggerEnabled = true
	}()

	mux.HandleFunc("/healthcheck", func(writer http.ResponseWriter, request *http.Request) {
		defer request.Body.Close()
		fmt.Fprint(writer, "ok")
	})

	//http-middlewares
	var middlewareChain middleware.ChainMiddleware
	if err = diContainer.Fill(middleware2.DefHttpMiddlewareChain, &middlewareChain); err != nil {
		return err
	}

	var logData = []zap.Field{
		zap.String("http_port", httpListen),
	}
	if swaggerEnabled {
		logData = append(
			logData,
			zap.Bool("swagger_enabled", swaggerEnabled),
			zap.String("swagger_path", swaggerPath),
		)
	}

	var m monitoring.Monitoring
	if err = diContainer.Fill(monitoring.DefMonitoring, &m); err != nil {
		return err
	}

	var rds redis.Redis
	if err = diContainer.Fill(redis_def.DefRedis, &rds); err != nil {
		return err
	}

	var counter counter2.Counter
	if err = diContainer.Fill(counter_def.DefCounter, &counter); err != nil {
		return err
	}

	var connCount int32
	var httpServer = http.Server{
		Addr:    httpListen,
		Handler: httpHandlerFunc(mux, middlewareChain),
		ConnState: func(conn net.Conn, state http.ConnState) {
			var delta int32
			switch state {
			case http.StateNew:
				delta = 1
			case http.StateClosed:
				delta = -1
			}
			atomic.AddInt32(&connCount, delta)
		},
	}

	go func() {
		var t = time.NewTicker(time.Second)
		for range t.C {
			if atomic.LoadInt32(&connCount) < 0 {
				continue
			}

			_ = m.Val(&monitoring.Metric{
				Namespace: "http",
				Name:      "connection_count",
			}, float64(connCount))
		}
	}()

	log.Info("start http server", logData...)

	go httpServer.ListenAndServe()
	select {
	case <-stopNotification:
		{
			err = counter.SaveDataFromIMDBToFile()
			return
		}
	}
}

func httpHandlerFunc(httpHandler http.Handler, chainMiddleware middleware.ChainMiddleware) http.Handler {
	return chainMiddleware.BuildChain(httpHandler.ServeHTTP)
}

func serveSwagger(mux *http.ServeMux) {
	mime.AddExtensionType(".svg", "image/svg+xml")
	var prefix = "/swagger-ui/"
	mux.Handle(prefix, http.StripPrefix(prefix, http.FileServer(&assetfs.AssetFS{
		Asset:    swagger.Asset,
		AssetDir: swagger.AssetDir,
	})))
}
