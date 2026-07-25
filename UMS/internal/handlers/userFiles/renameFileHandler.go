package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpcUtils "github.com/rakshithrajs/cloud/UMS/internal/grpc/utils"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

var fileRenamedSuccessfully = "file renamed successfully"

func (h *UserFilesHandler) RenameFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, fnRenameFile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if fileID == "" {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrFileIDRequired, fnRenameFile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	var payload models.RenameFileRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrInvalidJSON, fnRenameFile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	if err := modelUtils.Validate.Struct(&payload); err != nil {
		handlerUtils.ReturnErrorResponse(c, modelUtils.FieldErrors(err), fnRenameFile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	newName := strings.TrimSpace(payload.NewName)

	if err := h.client.RenameFileGrpcHandler(ctx, userID, fileID, newName); err != nil {
		status, msg := grpcUtils.MapGRPCError(err, handlerUtils.ErrFailedToRenameFile.Error())
		slog.Error(handlerUtils.LogPrefix(fnRenameFile)+"failed to rename file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
}
