// Package listeners provide dependency injection definitions.
package listeners

import (
	"google.golang.org/grpc"

	"sber_cloud/tw/healthcheck"
	"sber_cloud/tw/proto"

	"sber_cloud/tw/container"
)

// DefGRPCListenerHealthCheck definition name.
const DefGRPCListenerHealthCheck = "listener_grpc_health_check"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name:  DefGRPCListenerHealthCheck,
			Scope: DefGRPCPublicListenerScope,
			Build: func(container container.Container) (_ interface{}, err error) {
				return func(srv *grpc.Server) {
					counter.RegisterHealthServer(srv, healthcheck_v1.NewListener())
				}, nil
			},
		})
	})
}
