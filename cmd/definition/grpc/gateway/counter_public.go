package gateway

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	"sber_cloud/tw/container"
	"sber_cloud/tw/proto"
)

// DefGRPCUserPublicGateway definition name.
const DefGRPCUserPublicGateway = "gateway_grpc_user_public"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name:  DefGRPCUserPublicGateway,
			Scope: DefGRPCGatewayScope,
			Build: func(container container.Container) (_ interface{}, err error) {
				return func(mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error {
					return counter.RegisterCounterPublicHandlerFromEndpoint(
						context.Background(),
						mux,
						endpoint,
						opts,
					)
				}, nil
			},
		})
	})
}
