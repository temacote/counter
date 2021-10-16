package listeners

import (
	"google.golang.org/grpc"

	listener "sber_cloud/tw/cmd/grpc/listeners/public"
	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/logger"
	"sber_cloud/tw/proto"
)

// DefGRPCListenerUserPublic definition name.
const DefGRPCListenerUserPublic = "listener_grpc_user_public"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name:  DefGRPCListenerUserPublic,
			Scope: DefGRPCPublicListenerScope,
			Build: func(container container.Container) (_ interface{}, err error) {
				var (
					log logger.Logger
				)

				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				return func(srv *grpc.Server) {
					counter.RegisterCounterPublicServer(srv, listener.NewCounterPublicListener(
						log,
					))
				}, nil
			},
		})
	})
}
