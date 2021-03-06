package counter

import (
	"encoding/json"
	"fmt"
	"strings"
)

var (
	_ = fmt.Errorf
	_ = json.Unmarshal
	_ = strings.Replace
)

// Consul routings
func CounterPublicConsulRouting(serviceName string) map[string]string {
	type route struct {
		ServiceName string   `json:"service_name"`
		Name        string   `json:"name"`
		HttpMethods []string `json:"http_methods"`
		Route       string   `json:"route"`
		IsStream    bool     `json:"is_stream"`
	}

	var r = []*route{
		{
			ServiceName: serviceName,
			Name:        "CounterPublic.CountV1",
			HttpMethods: []string{"GET"},
			Route:       `/v1/count`,
			IsStream:    false,
		},
	}
	var (
		m    = map[string]string{}
		data []byte
		err  error
	)
	for _, c := range r {
		if data, err = json.Marshal(c); err != nil {
			panic(err)
		}

		var key = strings.Replace(fmt.Sprintf("route.%s", c.Name), ".", "_", -1)
		m[key] = string(data)
	}

	return m
}
