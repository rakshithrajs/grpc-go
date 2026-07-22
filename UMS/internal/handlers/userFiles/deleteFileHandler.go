package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"google.golang.org/grpc/metadata"
)

var deleteFileSuccessMsg = "file deleted successfully"

func (h *UserFilesHandler) DeleteFileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: err.Error()})
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if fileID == "" {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: handlers.ErrFileIDRequired.Error()})
		return
	}

	_, err = h.storage.GetUserFileName(c.Request.Context(), userID, fileID)
	if err != nil {
		slog.Error(handlers.LogPrefix(fnDeleteFile)+"failed to verify user file ownership", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}

	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "x-user-id", userID)

	if err := h.storage.DeleteUserFile(c.Request.Context(), userID, fileID); err != nil {
		slog.Error(handlers.LogPrefix(fnDeleteFile)+"failed to delete user file mapping", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}

	_, err = h.MMSClient.DeleteFile(ctx, &MMSpb.DeleteFileRequest{FileID: fileID})
	if err != nil {
		status, msg := handlers.MapGRPCError(err, handlers.ErrFailedToDeleteFile.Error())
		slog.Error(handlers.LogPrefix(fnDeleteFile)+"failed to delete file in MMS", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": deleteFileSuccessMsg})
}
