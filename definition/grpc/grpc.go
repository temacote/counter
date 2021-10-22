// Package monitoring provide dependency injection definitions.
package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"

	"sber_cloud/tw/cmd/definition/grpc/listeners"
	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/grpc/handlers"
	"sber_cloud/tw/definition/monitoring"
)

// DefGRPCServerBuilder definition name.
const DefGRPCServerBuilder = "grpc_server_builder"

type GrpcServerBuilder func(unaryMiddlewareDef, streamMiddlewareDef, listenerScope string) (*grpc.Server, error)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCServerBuilder,
			Build: func(container container.Container) (_ interface{}, err error) {
				return func(unaryMiddlewareDef, streamMiddlewareDef, listenerScope string) (*grpc.Server, error) {
					var (
						unaryServerOptions  = make([]grpc.UnaryServerInterceptor, 0, 8)
						streamServerOptions = make([]grpc.StreamServerInterceptor, 0, 8)
					)

					var unaryMiddlewareList []string
					if err = container.Fill(unaryMiddlewareDef, &unaryMiddlewareList); err != nil {
						return nil, err
					}

					var streamMiddlewareList []string
					if err = container.Fill(streamMiddlewareDef, &streamMiddlewareList); err != nil {
						return nil, err
					}

					for _, def := range unaryMiddlewareList {
						var opt []grpc.UnaryServerInterceptor
						if err = container.UnscopedFill(def, &opt); err != nil {
							return nil, err
						}
						unaryServerOptions = append(unaryServerOptions, opt...)
					}

					for _, def := range streamMiddlewareList {
						var opt []grpc.StreamServerInterceptor
						if err = container.UnscopedFill(def, &opt); err != nil {
							return nil, err
						}
						streamServerOptions = append(streamServerOptions, opt...)
					}

					var h handlers.StatHandler
					if err = container.Fill(handlers.DefGRPCMonitoringHandler, &h); err != nil {
						return nil, err
					}

					var grpcServer = grpc.NewServer([]grpc.ServerOption{
						grpc.StatsHandler(h),
						grpc_middleware.WithUnaryServerChain(unaryServerOptions...),
						grpc_middleware.WithStreamServerChain(streamServerOptions...),
					}...)

					for _, def := range container.Definitions() {
						if def.Scope == listenerScope {
							var registrant listeners.GRPCListenerRegistrant
							if err = container.UnscopedFill(def.Name, &registrant); err != nil {
								return nil, err
							}
							registrant(grpcServer)
						}
					}

					var m monitoring.Monitoring
					if err = container.Fill(monitoring.DefMonitoring, &m); err != nil {
						return nil, err
					}

					grpc_prometheus.EnableHandlingTimeHistogram()

					return grpcServer, nil
				}, nil
			},
		})
	})
}
