package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

var fileRenamedSuccessfully = "file renamed successfully"

func (h *UserFilesHandler) RenameFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	var payload models.RenameFileRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrInvalidJSON, FnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if err := modelUtils.Validate.Struct(&modelUtils.FileIDPayload{FileID: fileID}); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	if err := modelUtils.Validate.Struct(&payload); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	newName := strings.TrimSpace(payload.NewName)

	if err := h.client.RenameFileGrpcHandler(ctx, userID, fileID, newName); err != nil {
		slog.Error(handlerUtils.LogPrefix(FnRenameFile)+"failed to rename file", slog.Any(config.ErrorKey, err))
		status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrFailedToRenameFile.Error())
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
}
