package listeners

import "sber_cloud/tw/definition/logger"
import "sber_cloud/tw/proto"

type CounterPublicListener struct {
	log logger.Logger
}

func NewCounterPublicListener(
	log logger.Logger,
) counter.CounterPublicServer {
	return &CounterPublicListener{
		log: log,
	}
}
