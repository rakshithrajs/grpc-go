package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func (h *UserFilesHandler) ListFilesHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		utils.ReturnErrorResponse(c, err, fnListFiles, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	files, err := h.storage.ListUserFiles(ctx, userID)
	if err != nil {
		slog.Error(handlers.LogPrefix(fnListFiles)+"failed to list user files", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(c, err, fnListFiles, utils.ErrFailedToListFiles, "")
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
