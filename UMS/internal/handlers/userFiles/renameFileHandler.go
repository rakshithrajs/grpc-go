package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"google.golang.org/grpc/metadata"
)

var fileRenamedSuccessfully = "file renamed successfully"

func (h *UserFilesHandler) RenameFileHandler(c *gin.Context) {
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

	var payload models.RenameFileRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: handlers.ErrInvalidJSON.Error()})
		return
	}

	newName := strings.TrimSpace(payload.NewName)

	oldName, err := h.storage.GetUserFileName(c.Request.Context(), userID, fileID)
	if err != nil {
		slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to get user file name", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}
	if oldName == "" {
		c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
		return
	}

	if err := h.storage.UpdateUserFile(c.Request.Context(), userID, fileID, newName); err != nil {
		slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to update user file mapping", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}

	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "x-user-id", userID)
	_, err = h.MMSClient.RenameFile(ctx, &MMSpb.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	})
	if err != nil {
		if rbErr := h.storage.UpdateUserFile(c.Request.Context(), userID, fileID, oldName); rbErr != nil {
			slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to rollback user file mapping", slog.Any(config.ErrorKey, rbErr))
		}
		status, msg := handlers.MapGRPCError(err, handlers.ErrFailedToRenameFile.Error())
		slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to rename file in MMS", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
}
