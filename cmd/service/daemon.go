package service

import (
	"errors"
	"time"

	"go.uber.org/zap"

	"sber_cloud/tw/monitoring"
)

type daemon struct {
	monitoring monitoring.Monitoring
	logger     *zap.Logger
}

func NewDaemon(logger *zap.Logger, monitoring monitoring.Monitoring) Service {
	return &daemon{
		monitoring: monitoring,
		logger:     logger,
	}
}

func (d *daemon) Start() (err error) {
	var t = time.NewTicker(time.Second)
	for range t.C {
		d.logger.Error("hello", zap.Error(errors.New("some error")))

		if err = d.monitoring.Inc(&monitoring.Metric{
			Namespace: "dev",
			Subsystem: "counter",
			Name:      "daemon",
			ConstLabels: map[string]string{
				"method": "start",
			},
		}); err != nil {
			return err
		}
	}
	return nil
}
