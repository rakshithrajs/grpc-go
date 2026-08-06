package handlers

import (
	"net/http"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type grpcErrorMapping struct {
	status  int
	message string
}

var grpcCodeMapping = map[codes.Code]grpcErrorMapping{
	codes.AlreadyExists:   {status: http.StatusConflict},
	codes.Unauthenticated: {status: http.StatusUnauthorized, message: handlerErrors.ErrUnauthorized.Error()},
	codes.Internal:        {status: http.StatusInternalServerError},
}

// MapGRPCError maps a gRPC status error to an HTTP status code and message.
func MapGRPCError(err error, defaultMsg string) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, defaultMsg
	}

	mapping, ok := grpcCodeMapping[st.Code()]
	if !ok {
		return http.StatusInternalServerError, defaultMsg
	}

	if mapping.message != "" {
		return mapping.status, mapping.message
	}
	return mapping.status, st.Message()
}
