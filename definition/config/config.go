package config

import (
	"errors"
	"strings"

	"github.com/spf13/viper"

	"sber_cloud/tw/container"
)

// DefConfig definition name.
const DefConfig = "config"

type Config = *viper.Viper

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		var ok bool
		if _, ok = params["config"]; !ok {
			return errors.New("can't get required parameter config path")
		}

		var path string
		if path, ok = params["config"].(string); !ok {
			return errors.New(`parameter "config_path" should be string`)
		}

		return builder.Add(container.Def{
			Name: DefConfig,
			Build: func(container container.Container) (_ interface{}, err error) {
				var cfg = viper.New()

				cfg.AutomaticEnv()
				cfg.SetEnvPrefix("ENV")
				cfg.SetEnvKeyReplacer(
					strings.NewReplacer(".", "_"),
				)
				cfg.SetConfigFile(path)
				cfg.SetConfigType("yaml")

				if err = cfg.ReadInConfig(); err != nil {
					return nil, err
				}

				cfg.WatchConfig()
				cfg.SetDefault("monitoring.url", "http://localhost:9091")
				cfg.SetDefault("logger.url", "localhost:5110")
				cfg.SetDefault("consul.address", "localhost:8500")
				cfg.SetDefault("http.listen", "3018")
				cfg.SetDefault("grpc.public_listen", "127.0.0.1:3017")
				cfg.SetDefault("grpc.public_endpoint", "127.0.0.1:3017")
				cfg.SetDefault("tracing.url", "localhost:6831")
				cfg.SetDefault("redis.url", "localhost:6379")
				cfg.SetDefault("storage_file.path", "local.txt")

				return cfg, nil
			},
		})
	})
}
