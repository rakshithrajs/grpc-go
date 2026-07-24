package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

var deleteFileSuccessMsg = "file deleted successfully"

func (h *UserFilesHandler) DeleteFileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		utils.ReturnErrorResponse(c, err, fnDeleteFile, utils.ErrSomethingWentWrong, "")
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if fileID == "" {
		utils.ReturnErrorResponse(c, utils.ErrFileIDRequired, fnDeleteFile, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	if err := h.client.DeleteFileGrpcHandler(ctx, userID, fileID); err != nil {
		status, msg := handlers.MapGRPCError(err, utils.ErrFailedToDeleteFile.Error())
		slog.Error(handlers.LogPrefix(fnDeleteFile)+"failed to delete file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": deleteFileSuccessMsg})
}
