package gateway

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type (
	GRPCErrorHandler interface {
		HTTPError(context.Context, *runtime.ServeMux, runtime.Marshaler, http.ResponseWriter, *http.Request, error)
		HTTPStatusFromCode(code codes.Code) int
	}

	grpcErrorHandler struct {
		logger *zap.Logger
	}

	ErrorDetails struct {
		Error      string `json:"error"`
		Message    string `json:"message"`
		Translated string `json:"tMessage"`
	}

	ErrorWrapper struct {
		Body *ErrorBody `json:"exception"`
	}

	ErrorBody struct {
		Code       int32        `json:"code"`
		Error      ErrorDetails `json:"error,omitempty"`
		HTTPStatus int          `json:"httpStatus"`
	}
)

type (
	ValidationErrorMapUnit struct {
		Err        error  `json:"-"`
		Error      string `json:"error"`
		Message    string `json:"message"`
		Translated string `json:"tMessage"`
	}

	ValidationErrorMap struct {
		EmailError       []*ValidationErrorMapUnit `json:"email,omitempty"`
		PasswordError    []*ValidationErrorMapUnit `json:"password,omitempty"`
		PhoneError       []*ValidationErrorMapUnit `json:"phoneNumber,omitempty"`
		NameError        []*ValidationErrorMapUnit `json:"name,omitempty"`
		OldPasswordError []*ValidationErrorMapUnit `json:"oldPassword,omitempty"`
		SocialEmailError []*ValidationErrorMapUnit `json:"socialEmail, omitempty"`
		SocialTypeError  []*ValidationErrorMapUnit `json:"socialType, omitempty"`
		SocialIdError    []*ValidationErrorMapUnit `json:"socialId, omitempty"`
	}

	MapErrorWrapper struct {
		Body *ValidationMapErrorBody `json:"mapException"`
	}

	ValidationMapErrorBody struct {
		Code       int32              `json:"code"`
		Errors     ValidationErrorMap `json:"errors,omitempty"`
		HTTPStatus int                `json:"httpStatus"`
	}
)

type (
	ErrorTypeResolver struct {
		Exception    json.RawMessage
		MapException json.RawMessage
	}
)

func (err *MapErrorWrapper) Error() string {
	jsn, _ := json.Marshal(err)
	return string(jsn)
}

func (err *ErrorWrapper) Error() string {
	jsn, _ := json.Marshal(err)
	return string(jsn)
}

func NewValidationErrorMap() *ValidationErrorMap {
	return &ValidationErrorMap{
		EmailError:       []*ValidationErrorMapUnit{},
		PasswordError:    []*ValidationErrorMapUnit{},
		PhoneError:       []*ValidationErrorMapUnit{},
		NameError:        []*ValidationErrorMapUnit{},
		OldPasswordError: []*ValidationErrorMapUnit{},
		SocialEmailError: []*ValidationErrorMapUnit{},
		SocialTypeError:  []*ValidationErrorMapUnit{},
		SocialIdError:    []*ValidationErrorMapUnit{},
	}
}

func NewGRPCErrorHandler(logger *zap.Logger) GRPCErrorHandler {
	return &grpcErrorHandler{
		logger: logger,
	}
}

// кастомный маппер кодов от grpc на http
func (g *grpcErrorHandler) HTTPStatusFromCode(code codes.Code) int {
	switch code {
	//TODO implement error status switching
	}
	return 0
}

// HTTPError кастомный обработчик ошибок от grpc
func (g *grpcErrorHandler) HTTPError(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	_ *http.Request,
	err error,
) {
	const fallback = `{"exception": {"error": "failed to marshal error message"}}`

	w.Header().Del("Trailer")
	w.Header().Set("Content-Type", marshaler.ContentType())
	w.Header().Add("X-Request-Token", "")

	var (
		ok         bool
		s          *status.Status
		mapError   *MapErrorWrapper
		buff       []byte
		mErr       error
		body       interface{}
		statusCode int
	)

	if s, ok = status.FromError(err); !ok {
		s = status.New(codes.Unknown, err.Error())
	}

	if err = json.Unmarshal([]byte(s.Message()), &mapError); err == nil {
		body = mapError
		if statusCode = g.HTTPStatusFromCode(codes.Code(mapError.Body.Code)); statusCode == 0 {
			statusCode = runtime.HTTPStatusFromCode(s.Code())
		}
	} else {
		if statusCode = g.HTTPStatusFromCode(s.Code()); statusCode == 0 {
			statusCode = runtime.HTTPStatusFromCode(s.Code())
		}
		body = &ErrorWrapper{Body: &ErrorBody{
			Code: int32(s.Code()),
			Error: ErrorDetails{
				Message: s.Message(),
			},
		}}
	}

	switch body.(type) {
	case *MapErrorWrapper:
		body.(*MapErrorWrapper).Body.HTTPStatus = statusCode
	case *ErrorWrapper:
		body.(*ErrorWrapper).Body.HTTPStatus = statusCode
	}

	if buff, mErr = marshaler.Marshal(body); mErr != nil {
		g.logger.Error("Failed to marshal error message", zap.Reflect("body", body), zap.Error(mErr))
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := io.WriteString(w, fallback); err != nil {
			g.logger.Error("Failed to write response", zap.Error(err))
		}
		return
	}

	w.WriteHeader(statusCode)
	if _, err := w.Write(buff); err != nil {
		g.logger.Error("Failed to write response", zap.Error(err))
	}
}
