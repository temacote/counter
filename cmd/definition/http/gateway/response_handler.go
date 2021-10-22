// Package gateway provide dependency injection definitions.
package gateway

import (
	"sber_cloud/tw/cmd/http/gateway"
	"sber_cloud/tw/container"
	counter2 "sber_cloud/tw/definition/counter"
	"sber_cloud/tw/definition/logger"
)

const (
	// DefGRPCGatewayResponseHandler definition name.
	DefGRPCGatewayResponseHandler = "grpc_gateway_response_handler"
)

// GRPCErrorHandler type alias gateway.GRPCErrorHandler
type GRPCResponseHandler = gateway.GRPCResponseHandler

// Definition init func.
func init() {
	container.Register(func(builder *container.Builder, params map[string]interface{}) error {
		return builder.Add(container.Def{
			Name: DefGRPCGatewayResponseHandler,
			Build: func(container container.Container) (_ interface{}, err error) {
				var counter counter2.Counter
				if err = container.Fill(counter2.DefCounter, &counter); err != nil {
					return nil, err
				}

				_ = counter.LoadFromFileToIMDB()

				var log logger.Logger
				if err = container.Fill(logger.DefLogger, &log); err != nil {
					return nil, err
				}

				return gateway.NewGRPCResponseHandler(
					counter,
					log,
				), nil
			},
		})
	})
}
