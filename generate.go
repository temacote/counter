//go:generate protoc -I=. --go_out=plugins=grpc:. proto/healthcheck.proto
//go:generate protoc -I=. -I=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3/third_party/googleapis -I=$GOPATH/pkg/mod/github.com/grpc-ecosystem/grpc-gateway@v1.11.3 --go_out=plugins=grpc:. --consulroute_out=. --swagger_out=logtostderr=true:. --grpc-gateway_out=logtostderr=true:. proto/grpc_public.proto proto/common.proto
//go:generate gofmt -w ./proto/

package main

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	_ "github.com/grpc-ecosystem/grpc-gateway/utilities"
	_ "github.com/mailru/easyjson"
)
