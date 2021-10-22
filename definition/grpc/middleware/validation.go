// Package monitoring provide dependency injection definitions.
package middleware

import (
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/go-grpc-middleware/validator"

	"sber_cloud/tw/container"
)

const (
	// DefGRPCUnaryMiddlewareValidation definition name.
	DefGRPCUnaryMiddlewareValidation = "middleware_grpc_unary_validation"

	// DefGRPCStreamMiddlewareValidation definition name.
	DefGRPCStreamMiddlewareValidation = "middleware_grpc_stream_validation"
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewareValidation,
			Build: func(container container.Container) (_ interface{}, err error) {
				return []grpc.UnaryServerInterceptor{
					grpc_validator.UnaryServerInterceptor(),
				}, nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewareValidation,
			Build: func(container container.Container) (_ interface{}, err error) {
				return []grpc.StreamServerInterceptor{
					grpc_validator.StreamServerInterceptor(),
				}, nil
			},
		})
	})
}
