package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"

	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
)

func (h *UserFilesHandler) ListFilesHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, fnListFiles, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	files, err := h.storage.ListUserFiles(ctx, userID)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(fnListFiles)+"failed to list user files", slog.Any(config.ErrorKey, err))
		handlerUtils.ReturnErrorResponse(c, err, fnListFiles, handlerUtils.ErrFailedToListFiles, "")
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
