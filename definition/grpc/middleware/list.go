// Package monitoring provide dependency injection definitions.
package middleware

import (
	"sber_cloud/tw/container"
)

const (
	// DefGRPCUnaryMiddlewarePriorityListPublic definition name.
	DefGRPCUnaryMiddlewarePriorityListPublic = "middleware_grpc_unary_priority_list_public"

	// DefGRPCStreamMiddlewarePriorityListPublic definition name.
	DefGRPCStreamMiddlewarePriorityListPublic = "middleware_grpc_stream_priority_list_public"

	// DefGRPCUnaryMiddlewarePriorityListInternal definition name.
	DefGRPCUnaryMiddlewarePriorityListInternal = "middleware_grpc_unary_priority_list_internal"

	// DefGRPCStreamMiddlewarePriorityListInternal definition name.
	DefGRPCStreamMiddlewarePriorityListInternal = "middleware_grpc_stream_priority_list_internal"
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewarePriorityListPublic,
			Build: func(container container.Container) (_ interface{}, err error) {
				var list []string
				if err = container.Fill(DefGRPCUnaryMiddlewarePriorityListInternal, &list); err != nil {
					return nil, err
				}

				return append(list, []string{
					// TODO тут можно добавить мидлвару авторизации например
				}...), nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewarePriorityListPublic,
			Build: func(container container.Container) (_ interface{}, err error) {
				var list []string
				if err = container.Fill(DefGRPCStreamMiddlewarePriorityListInternal, &list); err != nil {
					return nil, err
				}

				return append(list, []string{
					// TODO тут можно добавить мидлвару авторизации например
				}...), nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCUnaryMiddlewarePriorityListInternal,
			Build: func(container container.Container) (_ interface{}, err error) {
				return []string{
					DefGRPCUnaryMiddlewareRecovery,
					DefGRPCUnaryMiddlewareTracing,
					DefGRPCUnaryMiddlewareRequestID,
					DefGRPCUnaryMiddlewareLogging,
					DefGRPCUnaryMiddlewareMonitoring,
					DefGRPCUnaryMiddlewareValidation,
				}, nil
			},
		})
	})

	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCStreamMiddlewarePriorityListInternal,
			Build: func(container container.Container) (_ interface{}, err error) {

				return []string{
					DefGRPCStreamMiddlewareRecovery,
					DefGRPCStreamMiddlewareTracing,
					DefGRPCStreamMiddlewareRequestID,
					DefGRPCStreamMiddlewareLogging,
					DefGRPCStreamMiddlewareMonitoring,
					DefGRPCStreamMiddlewareValidation,
				}, nil
			},
		})
	})
}
