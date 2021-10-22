// Package listeners provide dependency injection definitions.
package listeners

import "google.golang.org/grpc"

const (
	// DefGRPCPublicListenerScope scope name.
	DefGRPCPublicListenerScope = "scope_listener_public_grpc"
)

type GRPCListenerRegistrant func(srv *grpc.Server)
