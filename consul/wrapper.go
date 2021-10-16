//go:generate mockgen -package=consul -destination=wrapper_mock.go -source=wrapper.go
package consul

import (
	"sort"
	"sync"

	"github.com/hashicorp/consul/api"
)

type (
	Wrapper interface {
		Register(r *api.AgentServiceRegistration) error
		AliveServiceNameByTag(string) ([]string, error)
		Deregister() error
		Client() *api.Client
	}

	wrapper struct {
		client     *api.Client
		lock       *sync.Mutex
		registered []*api.AgentServiceRegistration
	}
)

func NewWrapper(client *api.Client) Wrapper {
	return &wrapper{
		client:     client,
		lock:       &sync.Mutex{},
		registered: make([]*api.AgentServiceRegistration, 0, 8),
	}
}

func (c *wrapper) AliveServiceNameByTag(tag string) (_ []string, err error) {
	var h api.HealthChecks
	if h, _, err = c.client.Health().State("passing", nil); err != nil {
		return nil, err
	}

	var raw = make([]string, 0, 8)
	for _, v := range h {
		if !containsString(v.ServiceTags, tag) {
			continue
		}

		raw = append(raw, v.ServiceName)
	}

	sort.Strings(raw)
	var out = make([]string, 0, 8)

	for i := 0; i < len(raw); i++ {
		if i == 0 || (i > 0 && raw[i] != raw[i-1]) {
			out = append(out, raw[i])
		}
	}

	return out, nil
}

func (c *wrapper) Register(r *api.AgentServiceRegistration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	c.registered = append(c.registered, r)
	return c.client.Agent().ServiceRegister(r)
}

func (c *wrapper) Deregister() (err error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for i := range c.registered {
		if err = c.client.Agent().ServiceDeregister(c.registered[i].ID); err != nil {
			return err
		}
	}

	return nil
}

func (c *wrapper) Client() *api.Client {
	return c.client
}

func containsString(s []string, v string) bool {
	for _, vv := range s {
		if vv == v {
			return true
		}
	}
	return false
}
