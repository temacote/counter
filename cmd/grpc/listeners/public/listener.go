package listeners

import (
	counter2 "sber_cloud/tw/definition/counter"
	"sber_cloud/tw/definition/logger"
	"sber_cloud/tw/proto"
)

type CounterPublicListener struct {
	log     logger.Logger
	counter counter2.Counter
}

func NewCounterPublicListener(
	log logger.Logger,
	counter counter2.Counter,
) counter.CounterPublicServer {
	return &CounterPublicListener{
		log:     log,
		counter: counter,
	}
}
