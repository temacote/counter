package counter

import (
	"sber_cloud/tw/container"
	"sber_cloud/tw/counter"
	"sber_cloud/tw/definition/config"
	redis_def "sber_cloud/tw/definition/redis"
	"sber_cloud/tw/redis"
)

// DefConsulWatcher definition name.
const DefCounter = "counter"

type Counter = counter.Counter

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefCounter,
			Build: func(container container.Container) (_ interface{}, err error) {
				var conf config.Config
				if err = container.Fill(config.DefConfig, &conf); err != nil {
					return nil, err
				}

				var rds redis.Redis
				if err = container.Fill(redis_def.DefRedis, &rds); err != nil {
					return nil, err
				}
				return counter.NewCounter(conf, rds), nil
			},
		})
	})
}
