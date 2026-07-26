package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
)

const (
	FuncNameAuthMiddleware = "AuthMiddleware"
	logPrefix              = "[" + FuncNameAuthMiddleware + "]: "
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader == config.NullString {
			handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrMissingAuthHeader, FuncNameAuthMiddleware, middlewareUtils.ErrSomethingWentWrong, config.NullString)
			c.Abort()
			return
		}

		claims, err := middlewareUtils.VerifyToken(authHeader)
		if err != nil {
			handlerErrors.ReturnErrorResponse(c, err, FuncNameAuthMiddleware, middlewareUtils.ErrSomethingWentWrong, config.NullString)
			c.Abort()
			return
		}

		userID := claims.Subject

		c.Set(config.UserIDMetadataKey, userID)
		c.Next()
	}
}
