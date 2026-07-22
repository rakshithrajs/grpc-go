package middleware

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

const (
	funcNameAuthMiddleware = "AuthMiddleware"
	logPrefix              = "[" + funcNameAuthMiddleware + "]: "
)

var (
	ErrMissingAuthHeader = errors.New("missing Authorization header")
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == config.NullString {
			c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: ErrMissingAuthHeader.Error()})
			c.Abort()
			return
		}

		claims, err := utils.VerifyToken(authHeader)
		if err != nil {
			if errors.Is(err, utils.ErrMissingBearerToken) || errors.Is(err, utils.ErrInvalidToken) || errors.Is(err, utils.ErrTokenExpired) {
				c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: err.Error()})
				c.Abort()
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: utils.ErrSomethingWentWrong.Error()})
			c.Abort()
			return
		}

		userID := claims.Subject

		c.Set("userID", userID)
		c.Next()
	}
}
