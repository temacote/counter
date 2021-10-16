// Package listeners provide dependency injection definitions.
package service

import (
	"sber_cloud/tw/cmd/service"
	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/logger"
	"sber_cloud/tw/definition/monitoring"
)

// DefServiceDaemon definition name.
const DefServiceDaemon = "service_daemon"

type Service = service.Service

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefServiceDaemon,
			Build: func(container container.Container) (_ interface{}, err error) {
				var m monitoring.Monitoring
				if err = container.Fill(monitoring.DefMonitoring, &m); err != nil {
					return nil, err
				}

				var l logger.Logger
				if err = container.Fill(logger.DefLogger, &l); err != nil {
					return nil, err
				}

				return service.NewDaemon(l, m), nil
			},
		})
	})
}
