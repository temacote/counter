// Package monitoring provide dependency injection definitions.
package middleware

import (
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/logger"
)

const (
	// DefGRPCUnaryMiddlewareLogging definition name.
	DefGRPCUnaryMiddlewareLogging = "middleware_grpc_unary_logging"

	// DefGRPCStreamMiddlewareLogging definition name.
	DefGRPCStreamMiddlewareLogging = "middleware_grpc_stream_logging"
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewareLogging,
			Build: func(container container.Container) (_ interface{}, err error) {
				var log logger.Logger
				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				return []grpc.UnaryServerInterceptor{
					grpc_zap.UnaryServerInterceptor(log),
				}, nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewareLogging,
			Build: func(container container.Container) (_ interface{}, err error) {
				var log logger.Logger
				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				return []grpc.StreamServerInterceptor{
					grpc_zap.StreamServerInterceptor(log),
				}, nil
			},
		})
	})
}
