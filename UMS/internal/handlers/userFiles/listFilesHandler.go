package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
)

func (h *UserFilesHandler) ListFilesHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnListFiles, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	files, err := h.storage.ListUserFiles(ctx, userID)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnListFiles)+"failed to list user files", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(c, err, FnListFiles, handlerErrors.ErrFailedToListFiles, config.NullString)
		return
	}

	c.JSON(http.StatusOK, gin.H{"files": files})
}
