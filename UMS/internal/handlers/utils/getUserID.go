package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
)

// GetUserIDFromGin retrieves the user ID from the Gin context. It returns an error if the user ID is not found.
func GetUserIDFromGin(c *gin.Context) (string, error) {
	userID, exists := c.Get(config.UserIDMetadataKey)
	if !exists {
		return config.NullString, handlerErrors.ErrUnauthorized
	}
	return userID.(string), nil
}
