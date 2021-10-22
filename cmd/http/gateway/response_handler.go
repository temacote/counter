package gateway

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/runtime/protoiface"

	"go.uber.org/zap"

	counter2 "sber_cloud/tw/counter"
	"sber_cloud/tw/definition/logger"
)

type (
	GRPCResponseHandler interface {
		ResponseHandler(ctx context.Context, w http.ResponseWriter, p protoiface.MessageV1) error
		MetadataHandler(ctx context.Context, rq *http.Request) metadata.MD
	}

	grpcResponseHandler struct {
		counter counter2.Counter
		logger  logger.Logger
	}
)

func NewGRPCResponseHandler(counter counter2.Counter, logger logger.Logger) GRPCResponseHandler {
	return &grpcResponseHandler{
		counter: counter,
		logger:  logger,
	}
}

// ResponseHandler обработчик ответа
func (h *grpcResponseHandler) ResponseHandler(ctx context.Context, w http.ResponseWriter, p protoiface.MessageV1) error {
	// TODO implement later
	return nil
}

// MetadataHandler обработчик формирвания метадаты ответа
func (h *grpcResponseHandler) MetadataHandler(ctx context.Context, rq *http.Request) metadata.MD {
	err := h.counter.AddToHistory(rq)
	if err != nil {
		h.logger.Error("error adding request to IMDB storage", zap.Error(err))
	}
	return metadata.MD{}
}
