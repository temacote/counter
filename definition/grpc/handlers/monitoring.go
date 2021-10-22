// Package monitoring provide dependency injection definitions.
package handlers

import (
	"google.golang.org/grpc/stats"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/monitoring"
	"sber_cloud/tw/grpc/handlers"
)

const DefGRPCMonitoringHandler = "grpc_handler_monitoring"

type StatHandler = stats.Handler

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCMonitoringHandler,
			Build: func(container container.Container) (_ interface{}, err error) {
				var m monitoring.Monitoring
				if err = container.Fill(monitoring.DefMonitoring, &m); err != nil {
					return nil, err
				}

				return handlers.NewMonitoringHandler(m), nil
			},
		})
	})

}
