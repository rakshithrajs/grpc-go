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
)

var deleteFileSuccessMsg = "file deleted successfully"

func (h *UserFilesHandler) DeleteFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, fnDeleteFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if fileID == config.NullString {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrFileIDRequired, fnDeleteFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	if err := h.client.DeleteFileGrpcHandler(ctx, userID, fileID); err != nil {
		status, msg := grpcUtils.MapGRPCError(err, handlerUtils.ErrFailedToDeleteFile.Error())
		slog.Error(handlerUtils.LogPrefix(fnDeleteFile)+"failed to delete file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": deleteFileSuccessMsg})
}
