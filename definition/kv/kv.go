// Package config provide dependency injection definitions.
package kv

import (
	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/config"
	"sber_cloud/tw/definition/consul"
	"sber_cloud/tw/definition/logger"
	kv "sber_cloud/tw/kv"
)

// DefConsulWatcher definition name.
const DefConsulWatcher = "config.consul_watcher"

type ConsulWatcher = kv.ConsulWatcher

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefConsulWatcher,
			Build: func(container container.Container) (_ interface{}, err error) {
				var log logger.Logger
				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				var c *consul.Client
				if err = container.Fill(consul.DefConsul, &c); err != nil {
					return nil, err
				}

				var conf config.Config
				if err = container.Fill(config.DefConfig, &conf); err != nil {
					return nil, err
				}

				return kv.NewConsulWatcher(c, log, conf), nil
			},
		})
	})
}
