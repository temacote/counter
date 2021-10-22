package middleware

import (
	"context"
	"sort"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	HeaderRequestID = "Request-Id"
	ContextKey      = "token"
)

func UnaryServerInterceptorRequestIDBuilder(skipFor []string) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (_ interface{}, err error) {
	sort.Strings(skipFor)

	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (_ interface{}, err error) {
		var i = sort.SearchStrings(skipFor, info.FullMethod[1:])
		if i < len(skipFor) && skipFor[i] == info.FullMethod[1:] {
			return handler(ctx, req)
		}

		if ctx, err = requestIDCtx(ctx); err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

func StreamServerInterceptorRequestIDBuilder(skipFor []string) func(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler,
) (err error) {
	sort.Strings(skipFor)

	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) (err error) {
		var i = sort.SearchStrings(skipFor, info.FullMethod[1:])
		if i < len(skipFor) && skipFor[i] == info.FullMethod[1:] {
			return handler(srv, ss)
		}

		var ctx context.Context
		if ctx, err = requestIDCtx(ss.Context()); err != nil {
			return err
		}

		var wrapped = grpc_middleware.WrapServerStream(ss)
		wrapped.WrappedContext = ctx

		return handler(srv, wrapped)
	}
}

func requestIDCtx(ctx context.Context) (context.Context, error) {
	var (
		md metadata.MD
		ok bool
	)

	// get metadata
	if md, ok = metadata.FromIncomingContext(ctx); !ok {
		return nil, status.Error(codes.InvalidArgument, "Metadata not sent")
	}

	// check request id exist
	var xRequestID string
	if len(md.Get(HeaderRequestID)) < 1 {
		return nil, status.Error(codes.InvalidArgument, HeaderRequestID+" not sent")
	}

	// check length
	xRequestID = md.Get(HeaderRequestID)[0]
	if len(xRequestID) == 0 {
		return nil, status.Error(codes.InvalidArgument, "Empty "+HeaderRequestID)
	}

	// if have tracing set request id to context
	if span := opentracing.SpanFromContext(ctx); span != nil {
		span.SetTag("request_id", xRequestID)
		ctx = opentracing.ContextWithSpan(ctx, span)
	}

	// return request id with context
	return context.WithValue(ctx, ContextKey, xRequestID), nil
}
