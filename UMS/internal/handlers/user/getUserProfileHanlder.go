package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

func (h *UMSHandler) GetUserProfileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: err.Error()})
		return
	}

	user, err := h.storage.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{config.ErrorKey: err.Error()})
			return
		}
		slog.Error(handlers.LogPrefix(fnGetUserProfile)+"failed to get user by id", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}

	user.Password = config.NullString

	c.JSON(http.StatusOK, gin.H{"user": user})
}
