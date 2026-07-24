package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/metadata"
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

	if _, err = h.storage.GetUserFileName(ctx, userID, fileID); err != nil {
		slog.Error(handlers.LogPrefix(fnDeleteFile)+"failed to verify user file ownership", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(c, err, fnDeleteFile, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "userID", userID)

	if err := h.storage.DeleteUserFile(ctx, userID, fileID); err != nil {
		slog.Error(handlers.LogPrefix(fnDeleteFile)+"failed to delete user file mapping", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(c, err, fnDeleteFile, utils.ErrSomethingWentWrong, "")
		return
	}

	if _, err = h.MMSClient.DeleteFile(ctx, &MMSpb.DeleteFileRequest{FileID: fileID}); err != nil {
		status, msg := handlers.MapGRPCError(err, utils.ErrFailedToDeleteFile.Error())
		slog.Error(handlers.LogPrefix(fnDeleteFile)+"failed to delete file in MMS", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": deleteFileSuccessMsg})
}
