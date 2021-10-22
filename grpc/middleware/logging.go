package middleware

import (
	"context"
	"path"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	"github.com/mailru/easyjson"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
)

// UnaryServerInterceptor returns a new unary server interceptors that adds zap.Logger to the context.
func UnaryLoggerInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		var (
			ok bool
			v  easyjson.Marshaler

			requestJson []byte
			err         error

			fields = serverCallFields(info.FullMethod)
		)

		if v, ok = req.(easyjson.Marshaler); ok {
			if requestJson, err = easyjson.Marshal(v); err == nil {
				fields = append(fields, zap.ByteString("grpc_request", requestJson))
			}
		}

		var newCtx = newLoggerForCall(ctx, logger, time.Now(), fields)

		return handler(newCtx, req)
	}
}

// StreamServerInterceptor returns a new streaming server interceptor for OpenTracing.
func StreamLoggerInterceptor(logger *zap.Logger) grpc.StreamServerInterceptor {
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		var (
			newCtx context.Context

			startTime = time.Now()
			fields    = serverCallFields(info.FullMethod)
		)

		newCtx = newLoggerForCall(stream.Context(), logger, startTime, fields)
		var wrappedStream = grpc_middleware.WrapServerStream(stream)
		wrappedStream.WrappedContext = newCtx

		return handler(srv, wrappedStream)
	}
}

func serverCallFields(fullMethodString string) []zapcore.Field {
	var service = path.Dir(fullMethodString)[1:]
	var method = path.Base(fullMethodString)

	return []zapcore.Field{
		zap.String("grpc.service", service),
		zap.String("grpc.method", method),
	}
}

func newLoggerForCall(ctx context.Context, logger *zap.Logger, start time.Time, fields []zapcore.Field) context.Context {
	var (
		deadline time.Time
		ok       bool
	)

	fields = append(fields, zap.String("grpc.start_time", start.Format(time.RFC3339)))
	if deadline, ok = ctx.Deadline(); ok {
		fields = append(fields, zap.String("grpc.request.deadline", deadline.Format(time.RFC3339)))
	}
	var callLog = logger.With(append(fields, ctxzap.TagsToFields(ctx)...)...)

	return ctxzap.ToContext(ctx, callLog)
}
