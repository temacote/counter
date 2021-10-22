// Package gateway provide dependency injection definitions.
package gateway

import (
	"sber_cloud/tw/cmd/http/gateway"
	"sber_cloud/tw/container"
	"sber_cloud/tw/definition/logger"
)

const (
	// DefGRPCGatewayErrorHandler definition name.
	DefGRPCGatewayErrorHandler = "grpc_gateway_error_handler"
)

// GRPCErrorHandler type alias gateway.GRPCErrorHandler
type GRPCErrorHandler = gateway.GRPCErrorHandler

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCGatewayErrorHandler,
			Build: func(container container.Container) (_ interface{}, err error) {
				var log logger.Logger
				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				return gateway.NewGRPCErrorHandler(log), nil
			},
		})
	})
}
