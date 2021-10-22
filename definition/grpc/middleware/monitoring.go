// Package monitoring provide dependency injection definitions.
package middleware

import (
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-prometheus"

	"sber_cloud/tw/container"
)

const (
	// DefGRPCUnaryMiddlewareMonitoring definition name.
	DefGRPCUnaryMiddlewareMonitoring = "middleware_grpc_unary_monitoring"

	// DefGRPCStreamMiddlewareMonitoring definition name.
	DefGRPCStreamMiddlewareMonitoring = "middleware_grpc_stream_monitoring"
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewareMonitoring,
			Build: func(container container.Container) (_ interface{}, err error) {
				return []grpc.UnaryServerInterceptor{
					grpc_prometheus.UnaryServerInterceptor,
				}, nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewareMonitoring,
			Build: func(container container.Container) (_ interface{}, err error) {
				return []grpc.StreamServerInterceptor{
					grpc_prometheus.StreamServerInterceptor,
				}, nil
			},
		})
	})
}
