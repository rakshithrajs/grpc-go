package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
)

const (
	funcNameAuthMiddleware = "AuthMiddleware"
	logPrefix              = "[" + funcNameAuthMiddleware + "]: "
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader == config.NullString {
			handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrMissingAuthHeader, funcNameAuthMiddleware, middlewareUtils.ErrSomethingWentWrong, "")
			c.Abort()
			return
		}

		claims, err := middlewareUtils.VerifyToken(authHeader)
		if err != nil {
			handlerUtils.ReturnErrorResponse(c, err, funcNameAuthMiddleware, middlewareUtils.ErrSomethingWentWrong, "")
			c.Abort()
			return
		}

		userID := claims.Subject

		c.Set(config.UserIDMetadataKey, userID)
		c.Next()
	}
}
