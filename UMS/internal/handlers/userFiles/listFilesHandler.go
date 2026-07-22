package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
)

func (h *UserFilesHandler) ListFilesHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: err.Error()})
		return
	}

	files, err := h.storage.ListUserFiles(c.Request.Context(), userID)
	if err != nil {
		slog.Error(handlers.LogPrefix(fnListFiles)+"failed to list user files", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrFailedToListFiles.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
