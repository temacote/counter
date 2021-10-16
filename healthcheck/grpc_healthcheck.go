package healthcheck_v1

import (
	"context"

	"sber_cloud/tw/proto"
)

type HealthCheck struct{}

func NewListener() counter.HealthServer {
	return &HealthCheck{}
}

func (*HealthCheck) Check(context.Context, *counter.HealthCheckRequest) (*counter.HealthCheckResponse, error) {
	return &counter.HealthCheckResponse{
		Status: counter.HealthCheckResponse_SERVING,
	}, nil
}
