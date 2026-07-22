package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var LogPrefix = func(fnName string) string {
	return "[" + fnName + "]"
}

func GetUserIDFromGin(c *gin.Context) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", ErrUnauthorized
	}
	return userID.(string), nil
}

func MapGRPCError(err error, defaultMsg string) (int, string) {
	st, ok := status.FromError(err)
	if !ok {
		return http.StatusInternalServerError, defaultMsg
	}

	switch st.Code() {
	case codes.InvalidArgument:
		return http.StatusBadRequest, st.Message()
	case codes.AlreadyExists:
		return http.StatusConflict, st.Message()
	case codes.Unauthenticated:
		return http.StatusUnauthorized, ErrUnauthorized.Error()
	default:
		return http.StatusInternalServerError, defaultMsg
	}
}

func HandleDomainError(c *gin.Context, err error) bool {
	status, ok := domainErrorStatus(err)
	if !ok {
		return false
	}
	c.JSON(status, gin.H{config.ErrorKey: err.Error()})
	return true
}

func domainErrorStatus(err error) (int, bool) {
	switch {
	case errors.Is(err, storage.ErrUserEmailAlreadyExists),
		errors.Is(err, storage.ErrPhoneNumberAlreadyExists):
		return http.StatusConflict, true
	}
	return 0, false
}
