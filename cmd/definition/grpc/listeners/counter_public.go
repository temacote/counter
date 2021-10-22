package listeners

import (
	"google.golang.org/grpc"

	listener "sber_cloud/tw/cmd/grpc/listeners/public"
	"sber_cloud/tw/container"
	counter2 "sber_cloud/tw/definition/counter"
	"sber_cloud/tw/definition/logger"
	counter "sber_cloud/tw/proto"
)

// DefGRPCListenerUserPublic definition name.
const DefGRPCListenerCounterPublic = "listener_grpc_counter_public"

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name:  DefGRPCListenerCounterPublic,
			Scope: DefGRPCPublicListenerScope,
			Build: func(container container.Container) (_ interface{}, err error) {
				var (
					log logger.Logger
					cnt counter2.Counter
				)

				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				if err = container.Fill(counter2.DefCounter, &cnt); err != nil {
					return nil, err
				}

				return func(srv *grpc.Server) {
					counter.RegisterCounterPublicServer(srv, listener.NewCounterPublicListener(
						log,
						cnt,
					))
				}, nil
			},
		})
	})
}
