// Package redis provide dependency injection definitions.
package redis

import (
	"github.com/go-redis/redis"

	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/config"

	redis2 "sber_cloud/tw/redis"
)

// DefRedis definition name
const DefRedis = "redis"

// Pool type alias
type RedisClient = redis.Client

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefRedis,
			Build: func(container container.Container) (_ interface{}, err error) {
				var conf config.Config
				if err = container.Fill(config.DefConfig, &conf); err != nil {
					return nil, err
				}
				return redis2.NewRedis(conf), nil
			},
		})
	})
}
