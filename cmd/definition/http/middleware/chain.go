package middleware

import (
	"sber_cloud/tw/cmd/http/middleware"
	"sber_cloud/tw/container"
)

const (
	DefHttpMiddlewareChain = "http_middleware"
)

var middlewareList = []string{
	//Сюда добавляются Def-ключи нужных мидлвар
}

func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefHttpMiddlewareChain,
			Build: func(container container.Container) (_ interface{}, err error) {

				var c []middleware.HttpMiddleware

				for _, i := range middlewareList {
					var d middleware.HttpMiddleware
					if err = container.Fill(i, &d); err != nil {
						return nil, err
					}

					c = append(c, d)
				}

				return middleware.NewChainMiddleware(c...), nil
			},
		})
	})
}
