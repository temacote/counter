// Package monitoring provide dependency injection definitions.
package monitoring

import (
	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/config"
	"sber_cloud/tw/definition/logger"
	"sber_cloud/tw/monitoring"
)

// DefMonitoring definition name.
const DefMonitoring = "monitoring"

type (
	Monitoring = monitoring.Monitoring
	Metric     = monitoring.Metric

	monitoringConfig struct {
		Url      string `yaml:"url"`
		Username string `yaml:"username"`
		Password string `yaml:"password"`
	}
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		var subProcess = params["sub_process"].(string)
		return builder.Add(container.Def{
			Name: DefMonitoring,
			Build: func(container container.Container) (_ interface{}, err error) {
				var cfg config.Config
				if err = container.Fill(config.DefConfig, &cfg); err != nil {
					return nil, err
				}

				var l logger.Logger
				if err = container.Fill(logger.DefLogger, &l); err != nil {
					return nil, err
				}

				var conf = &monitoringConfig{}
				if err = cfg.UnmarshalKey("monitoring", conf); err != nil {
					return nil, err
				}

				return monitoring.NewPrometheusMonitoring(
					l,
					conf.Url,
					conf.Username,
					conf.Password,
					cfg.GetString("service"),
					subProcess,
					cfg.GetBool("monitoring.disable_log_push_error"),
				), nil
			},
		})
	})
}
