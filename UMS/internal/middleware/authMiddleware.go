package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

const (
	funcNameAuthMiddleware = "AuthMiddleware"
	logPrefix              = "[" + funcNameAuthMiddleware + "]: "
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := strings.TrimSpace(c.GetHeader("Authorization"))

		if authHeader == config.NullString {
			utils.ReturnErrorResponse(c, utils.ErrMissingAuthHeader, funcNameAuthMiddleware, utils.ErrSomethingWentWrong, "")
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(authHeader)
		if err != nil {
			utils.ReturnErrorResponse(c, err, funcNameAuthMiddleware, utils.ErrSomethingWentWrong, "")
			c.Abort()
			return
		}

		userID := claims.Subject

		c.Set(config.UserIDMetadataKey, userID)
		c.Next()
	}
}
