package customerrors

import "google.golang.org/grpc/status"

const (
	LoggerKeyToken = "token"
)

var (
	InternalServerError = status.Error(11501, "internal server error")
)
