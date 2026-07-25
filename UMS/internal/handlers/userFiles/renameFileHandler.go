package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

var fileRenamedSuccessfully = "file renamed successfully"

func (h *UserFilesHandler) RenameFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, fnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	var payload models.RenameFileRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrInvalidJSON, fnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	var fileID = strings.TrimSpace(c.Param("fileID"))
	if err := modelUtils.Validate.Var(fileID, "required,isValueEmpty,uuid"); err != nil {
		handlerUtils.ReturnErrorResponse(c, modelUtils.ErrFileIDRequired, fnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	if err := modelUtils.Validate.Struct(&payload); err != nil {
		handlerUtils.ReturnErrorResponse(c, modelUtils.FieldErrors(err), fnRenameFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	newName := strings.TrimSpace(payload.NewName)

	status, msg := h.client.RenameFileGrpcHandler(ctx, userID, fileID, newName)
	if status != http.StatusOK {
		slog.Error(handlerUtils.LogPrefix(fnRenameFile)+"failed to rename file", slog.Any(config.ErrorKey, msg))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
}
