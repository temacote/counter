package cmd

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"sber_cloud/tw/cmd/definition/grpc/listeners"
	"sber_cloud/tw/consul"
	consulDef "sber_cloud/tw/definition/consul"
	grpc_def "sber_cloud/tw/definition/grpc"
	"sber_cloud/tw/definition/grpc/middleware"
	"sber_cloud/tw/definition/logger"
)

var (
	grpcCmdPublic = &cobra.Command{
		Use:   "grpc_public",
		Short: "Run public GRPC server",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var (
				g      *grpc.Server
				listen = conf.GetString("grpc.public_listen")
			)

			var consulWrapper consulDef.Wrapper
			if err = diContainer.Fill(consulDef.DefConsulWrapper, &consulWrapper); err != nil {
				return err
			}

			if g, err = builder(
				listen,
				middleware.DefGRPCUnaryMiddlewarePriorityListPublic,
				middleware.DefGRPCStreamMiddlewarePriorityListPublic,
				listeners.DefGRPCPublicListenerScope,
			); err != nil {
				return err
			}

			var (
				strpPort string
				port     int64
			)
			if _, strpPort, err = net.SplitHostPort(listen); err != nil {
				return errors.Wrap(err, "error get listen grpc port from string")
			}

			if port, err = strconv.ParseInt(strpPort, 10, 64); err != nil {
				return errors.Wrap(err, "error parse grpc port from string")
			}

			if err = consulRegister(consulWrapper, "grpc_public", port, g); err != nil {
				return err
			}

			select {
			case <-stopNotification:
				{
					return consulWrapper.Deregister()
				}
			}
		},
	}
)

// Command init function.
func init() {
	rootCmd.AddCommand(grpcCmdPublic)
}

func builder(listen, uMidDef, sMidDef, listScope string) (g *grpc.Server, err error) {
	var log logger.Logger
	if err = diContainer.Fill(logger.DefLogger, &log); err != nil {
		return nil, err
	}

	var gb grpc_def.GrpcServerBuilder
	if err = diContainer.Fill(grpc_def.DefGRPCServerBuilder, &gb); err != nil {
		return nil, err
	}

	if g, err = gb(uMidDef, sMidDef, listScope); err != nil {
		return nil, err
	}

	grpc_prometheus.Register(g)

	var listener net.Listener
	if listener, err = net.Listen("tcp", listen); err != nil {
		return nil, err
	}

	log.Info(
		"Start GRPC server",
		zap.String("listen", listen),
	)
	go g.Serve(listener)

	return g, nil

}

func consulRegister(consulWrapper consul.Wrapper, tag string, port int64, g *grpc.Server) (err error) {
	var (
		serviceName       = conf.GetString("service")
		instanceName      = conf.GetString("instance")
		consulServiceName = serviceName + "-" + tag
	)

	type route struct {
		ServiceName string `json:"service_name"`
		Route       string `json:"route"`
		IsGRPC      bool   `json:"is_grpc"`
	}

	var grpcServiceList = make(map[string]string)
	for key := range g.GetServiceInfo() {
		if key == "grpc.health.v1.Health" {
			continue
		}

		var routeEncoded []byte
		if routeEncoded, err = json.Marshal(&route{
			ServiceName: consulServiceName,
			IsGRPC:      true,
			Route:       key,
		}); err != nil {
			return errors.Wrap(err, "error encode consul metadata")
		}

		key = strings.Replace(key, ".", "_", -1)

		grpcServiceList["route_"+key] = string(routeEncoded)
	}

	if err = consulWrapper.Register(&consulDef.AgentServiceRegistration{
		ID:   fmt.Sprintf("%s:grpc:%s:%s", serviceName, tag, instanceName),
		Name: consulServiceName,
		Tags: []string{"grpc", tag},
		Meta: grpcServiceList,
		Check: &consulDef.AgentServiceCheck{
			GRPC:     fmt.Sprintf("localhost:%d", port),
			Interval: "3s",
			Timeout:  "1s",
		},
		Port: int(port),
	}); err != nil {
		return errors.Wrap(err, "error register grpc service in consul")
	}

	return nil
}
