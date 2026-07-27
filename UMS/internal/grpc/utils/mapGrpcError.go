package grpc

import (
	"net/http"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func MapGRPCError(err error, defaultMsg string) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, defaultMsg
	}

	switch st.Code() {
	case codes.AlreadyExists:
		return http.StatusConflict, st.Message()
	case codes.Unauthenticated:
		return http.StatusUnauthorized, handlerErrors.ErrUnauthorized.Error()
	default:
		return http.StatusInternalServerError, defaultMsg
	}
}
