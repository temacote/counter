// Package consul provide dependency injection definitions.
package consul

import (
	consulApi "github.com/hashicorp/consul/api"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/config"
)

const (
	// DefConsul definition name.
	DefConsul = "consul"
)

type (
	// Client type alias of consul.Client.
	Client = consulApi.Client
	// AgentServiceRegistration type alias of consul.AgentServiceRegistration.
	AgentServiceRegistration = consulApi.AgentServiceRegistration
	// AgentServiceCheck type alias of consul.AgentServiceCheck.
	AgentServiceCheck = consulApi.AgentServiceCheck
)

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefConsul,
			Build: func(container container.Container) (_ interface{}, err error) {
				var cfg config.Config
				if err = container.Fill(config.DefConfig, &cfg); err != nil {
					return nil, err
				}

				var consulConf = consulApi.DefaultConfig()
				consulConf.Address = cfg.GetString("consul.address")
				consulConf.Token = cfg.GetString("consul.token")

				var client *Client
				if client, err = consulApi.NewClient(consulConf); err != nil {
					return nil, err
				}

				return client, nil
			},
		})
	})
}
