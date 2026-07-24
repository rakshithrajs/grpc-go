package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var LogPrefix = func(fnName string) string {
	return "[" + fnName + "]"
}

func GetUserIDFromGin(c *gin.Context) (string, error) {
	userID, exists := c.Get(config.UserIDMetadataKey)
	if !exists {
		return "", utils.ErrUnauthorized
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
		return http.StatusUnauthorized, utils.ErrUnauthorized.Error()
	default:
		return http.StatusInternalServerError, defaultMsg
	}
}
