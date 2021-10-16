// Package gateway provide dependency injection definitions.
package gateway

import (
	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

// DefGRPCGatewayScope scope name.
const DefGRPCGatewayScope = "gateway_scope_grpc"

type GRPCGatewayRegistrant func(*runtime.ServeMux, string, []grpc.DialOption) error
