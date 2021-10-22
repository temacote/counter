package gateway

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/config"
	consulDef "sber_cloud/tw/definition/consul"
	counter "sber_cloud/tw/proto"

	grpcGateway "sber_cloud/tw/cmd/definition/grpc/gateway"
	httpGateway "sber_cloud/tw/cmd/http/gateway"
)

// DefHTTPGateway definition name.
const DefHTTPGateway = "gateway_http"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefHTTPGateway,
			Build: func(container container.Container) (_ interface{}, err error) {
				var conf config.Config
				if err = container.Fill(config.DefConfig, &conf); err != nil {
					return nil, err
				}

				var errorHandler GRPCErrorHandler
				if err = container.Fill(DefGRPCGatewayErrorHandler, &errorHandler); err != nil {
					return nil, err
				}

				var responseHandler httpGateway.GRPCResponseHandler
				if err = container.Fill(DefGRPCGatewayResponseHandler, &responseHandler); err != nil {
					return nil, err
				}

				var consulWrapper consulDef.Wrapper
				if err = container.Fill(consulDef.DefConsulWrapper, &consulWrapper); err != nil {
					return nil, err
				}

				var (
					serviceName    = conf.GetString("service")
					instanceName   = conf.GetString("instance")
					grpcEndpoint   = conf.GetString("grpc.public_endpoint")
					httpBindString = conf.GetString("http.listen")
					httpBind       = strings.Split(httpBindString, ":")
				)

				if len(httpBind) != 2 {
					return nil, errors.New("incorrect http.listen value. Required host:port format")
				}

				var port int
				if port, err = strconv.Atoi(httpBind[1]); err != nil {
					return nil, errors.Wrap(err, "error parse http port on config http.listen")
				}

				var healthCheckAddr = httpBind[0]
				if len(healthCheckAddr) == 0 {
					healthCheckAddr = "localhost"
				}

				if err = consulWrapper.Register(&consulDef.AgentServiceRegistration{
					ID:   fmt.Sprintf("%s:http:%s", serviceName, instanceName),
					Name: serviceName + "-http",
					Tags: []string{"http"},
					Meta: counter.CounterPublicConsulRouting(serviceName + "-http"),
					Check: &consulDef.AgentServiceCheck{
						HTTP:     fmt.Sprintf("http://%s/healthcheck", net.JoinHostPort(healthCheckAddr, httpBind[1])),
						Interval: "3s",
						Timeout:  "1s",
					},
					Port: port,
				}); err != nil {
					return nil, errors.Wrap(err, "error register service in consul")
				}

				var (
					headerMapper = func(key string) (string, bool) {
						const keyPrefix = "X-"
						if strings.HasPrefix(key, keyPrefix) {
							return key[len(keyPrefix):], true
						}
						return runtime.DefaultHeaderMatcher(key)
					}

					mux = http.NewServeMux()
					// кастомный маппинг хедерров
					gwMux = runtime.NewServeMux(
						runtime.WithMarshalerOption("application/json", &runtime.JSONBuiltin{}),
						runtime.WithMarshalerOption("*", &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
						runtime.WithIncomingHeaderMatcher(headerMapper),
						runtime.WithForwardResponseOption(responseHandler.ResponseHandler),
						runtime.WithMetadata(responseHandler.MetadataHandler),
					)
				)

				// кастомный маппер ошибок
				runtime.HTTPError = errorHandler.HTTPError

				for _, def := range container.Definitions() {
					if def.Scope == grpcGateway.DefGRPCGatewayScope {
						var m grpcGateway.GRPCGatewayRegistrant
						if err = container.UnscopedFill(def.Name, &m); err != nil {
							return nil, err
						}
						if err = m(gwMux, grpcEndpoint, []grpc.DialOption{grpc.WithInsecure()}); err != nil {
							return nil, err
						}
					}
				}

				mux.Handle("/", gwMux)
				return mux, nil
			},
		})
	})
}
