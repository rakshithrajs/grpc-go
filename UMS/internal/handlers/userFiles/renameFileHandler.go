package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

var fileRenamedSuccessfully = "file renamed successfully"

// RenameFileHandler renames a user file via the MMS client.
func (h *UserFilesHandler) RenameFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnRenameFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	ctx := c.Request.Context()

	var payload models.RenameFileRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrInvalidJSON, FnRenameFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	var uri models.FileIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrInvalidURI, FnRenameFile, handlerErrors.ErrSomethingWentWrong)
		return
	}
	if err := modelUtils.Validate.Struct(&uri); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnRenameFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	if err := modelUtils.Validate.Struct(&payload); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnRenameFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	newName := strings.TrimSpace(payload.NewName)

	if err := h.client.RenameFileGrpcClient(ctx, userID, uri.FileID, newName); err != nil {
		slog.Error(handlerUtils.LogPrefix(FnRenameFile)+"failed to rename file", slog.Any(config.ErrorKey, err))
		status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrFailedToRenameFile.Error())
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
}
