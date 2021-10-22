// Package monitoring provide dependency injection definitions.
package middleware

import (
	"google.golang.org/grpc"

	"sber_cloud/tw/container"
	"sber_cloud/tw/grpc/middleware"
)

const (
	// DefGRPCUnaryMiddlewareRequestID definition name.
	DefGRPCUnaryMiddlewareRequestID = "middleware_grpc_unary_request_id"

	// DefGRPCStreamMiddlewareRequestID definition name.
	DefGRPCStreamMiddlewareRequestID = "middleware_grpc_stream_request_id"

	healthCheckRoute = "grpc.health.v1.Health/Check"
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewareRequestID,
			Build: func(container container.Container) (_ interface{}, err error) {
				return []grpc.UnaryServerInterceptor{
					middleware.UnaryServerInterceptorRequestIDBuilder([]string{
						healthCheckRoute,
					}),
				}, nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewareRequestID,
			Build: func(container container.Container) (_ interface{}, err error) {
				return []grpc.StreamServerInterceptor{
					middleware.StreamServerInterceptorRequestIDBuilder([]string{
						healthCheckRoute,
					}),
				}, nil
			},
		})
	})
}
