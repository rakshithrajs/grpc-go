package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	tmsGrpc "github.com/rakshithrajs/cloud/UMS/internal/grpcClient"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
)

const (
	// FuncNameAuthMiddleware is the log prefix for the auth middleware.
	FuncNameAuthMiddleware = "AuthMiddleware"
)

// AuthMiddleware validates the Authorization header and attaches the user ID to the Gin context.
func AuthMiddleware(tmsClient *tmsGrpc.TMSClient) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader == config.NullString {
			handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrMissingAuthHeader, FuncNameAuthMiddleware, handlerErrors.ErrSomethingWentWrong)
			c.Abort()
			return
		}

		claims, err := tmsClient.ValidateToken(c.Request.Context(), authHeader)
		if err != nil {
			status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrSomethingWentWrong.Error())
			c.JSON(status, gin.H{config.ErrorKey: errMsg})
			c.Abort()
			return
		}

		c.Set(config.UserIDMetadataKey, claims.UserID)
		c.Next()
	}
}
