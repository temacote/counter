// Package monitoring provide dependency injection definitions.
package middleware

import (
	"context"
	"path"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/logger"
	"sber_cloud/tw/definition/monitoring"
)

const (
	// DefGRPCUnaryMiddlewareRecovery definition name.
	DefGRPCUnaryMiddlewareRecovery = "middleware_grpc_unary_recovery"

	// DefGRPCStreamMiddlewareRecovery definition name.
	DefGRPCStreamMiddlewareRecovery = "middleware_grpc_stream_recovery"
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewareRecovery,
			Build: func(container container.Container) (_ interface{}, err error) {
				var log logger.Logger
				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				var m monitoring.Monitoring
				if err = container.Fill(monitoring.DefMonitoring, &m); err != nil {
					return nil, err
				}

				return []grpc.UnaryServerInterceptor{
					func(
						ctx context.Context,
						req interface{},
						info *grpc.UnaryServerInfo,
						handler grpc.UnaryHandler,
					) (_ interface{}, err error) {
						defer func() {
							if r := recover(); r != nil {
								var service, method = serverInfo(info.FullMethod)

								_ = m.Inc(&monitoring.Metric{
									Namespace: "grpc",
									Name:      "panic",
									ConstLabels: map[string]string{
										"grpc_service": service,
										"grpc_method":  method,
									},
								})

								log.Error(
									"grpc panic",
									zap.String("service", service),
									zap.String("method", method),
									zap.Reflect("data", r),
								)
								err = status.New(codes.Internal, "internal error").Err()
							}
						}()

						return handler(ctx, req)
					},
				}, nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewareRecovery,
			Build: func(container container.Container) (_ interface{}, err error) {
				var log logger.Logger
				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				var m monitoring.Monitoring
				if err = container.Fill(monitoring.DefMonitoring, &m); err != nil {
					return nil, err
				}

				return []grpc.StreamServerInterceptor{
					func(
						srv interface{},
						stream grpc.ServerStream,
						info *grpc.StreamServerInfo,
						handler grpc.StreamHandler,
					) (err error) {
						defer func() {
							if r := recover(); r != nil {
								var service, method = serverInfo(info.FullMethod)

								_ = m.Inc(&monitoring.Metric{
									Namespace: "grpc",
									Name:      "panic",
									ConstLabels: map[string]string{
										"grpc_service": service,
										"grpc_method":  method,
									},
								})

								log.Error(
									"grpc panic",
									zap.String("service", service),
									zap.String("method", method),
									zap.Reflect("data", r),
								)
								err = status.New(codes.Internal, "internal error").Err()
							}
						}()

						return handler(srv, stream)
					},
				}, nil
			},
		})
	})
}

func serverInfo(fullMethodString string) (service, method string) {
	return path.Dir(fullMethodString)[1:], path.Base(fullMethodString)
}
