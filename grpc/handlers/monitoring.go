package handlers

import (
	"context"
	"sync/atomic"

	"google.golang.org/grpc/stats"

	"sber_cloud/tw/monitoring"
)

type monitoringHandler struct {
	m monitoring.Monitoring
	c int32
}

func NewMonitoringHandler(m monitoring.Monitoring) stats.Handler {
	return &monitoringHandler{
		m: m,
	}
}

// HandleConn exists to satisfy gRPC stats.Handler.
func (s *monitoringHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
	var delta int32
	switch cs.(type) {
	case *stats.ConnEnd:
		delta = -1
	case *stats.ConnBegin:
		delta = 1
	}

	atomic.AddInt32(&s.c, delta)

	var val float64
	if s.c > 0 {
		val = float64(s.c)
	}

	_ = s.m.Val(&monitoring.Metric{
		Namespace: "grpc",
		Name:      "connection_count",
	}, val)
}

// TagConn exists to satisfy gRPC stats.Handler.
func (s *monitoringHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC tracing and stats instrumentation.
func (s *monitoringHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	// no-op
}

// TagRPC implements per-RPC context management.
func (s *monitoringHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	// no-op
	return ctx
}
