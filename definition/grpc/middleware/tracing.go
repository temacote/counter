// Package monitoring provide dependency injection definitions.
package middleware

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	"google.golang.org/grpc"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/tracing"
)

const (
	// DefGRPCUnaryMiddlewareTracing definition name.
	DefGRPCUnaryMiddlewareTracing = "middleware_grpc_unary_tracing"

	// DefGRPCStreamMiddlewareTracing definition name.
	DefGRPCStreamMiddlewareTracing = "middleware_grpc_stream_tracing"
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewareTracing,
			Build: func(container container.Container) (_ interface{}, err error) {
				var t tracing.Tracer
				if err = container.Fill(tracing.DefTracing, &t); err != nil {
					return nil, err
				}

				return []grpc.UnaryServerInterceptor{
					grpc_opentracing.UnaryServerInterceptor(
						grpc_opentracing.WithFilterFunc(func(ctx context.Context, fullMethodName string) bool {
							switch fullMethodName {
							case "/" + healthCheckRoute:
								return false
							default:
								return true
							}
						}),
						grpc_opentracing.WithTracer(t),
					),
				}, nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewareTracing,
			Build: func(container container.Container) (_ interface{}, err error) {
				var t tracing.Tracer
				if err = container.Fill(tracing.DefTracing, &t); err != nil {
					return nil, err
				}

				return []grpc.StreamServerInterceptor{
					grpc_opentracing.StreamServerInterceptor(
						grpc_opentracing.WithFilterFunc(func(ctx context.Context, fullMethodName string) bool {
							switch fullMethodName {
							case "/" + healthCheckRoute:
								return false
							default:
								return true
							}
						}),
						grpc_opentracing.WithTracer(t),
					),
				}, nil
			},
		})
	})
}
