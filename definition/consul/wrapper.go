// Package consul provide dependency injection definitions.
package consul

import (
	"sber_cloud/tw/consul"
	"sber_cloud/tw/container"
)

// DefConsulWrapper definition name.
const DefConsulWrapper = "consul_wrapper"

type Wrapper = consul.Wrapper

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefConsulWrapper,
			Build: func(container container.Container) (_ interface{}, err error) {
				var c *Client
				if err = container.Fill(DefConsul, &c); err != nil {
					return nil, err
				}

				return consul.NewWrapper(c), nil
			},
			Close: func(obj interface{}) error {
				return obj.(consul.Wrapper).Deregister()
			},
		})
	})
}
