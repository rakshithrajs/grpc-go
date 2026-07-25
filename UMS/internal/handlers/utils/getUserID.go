package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
)

func GetUserIDFromGin(c *gin.Context) (string, error) {
	userID, exists := c.Get(config.UserIDMetadataKey)
	if !exists {
		return "", ErrUnauthorized
	}
	return userID.(string), nil
}
