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
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/metadata"
)

var fileRenamedSuccessfully = "file renamed successfully"

func (h *UserFilesHandler) RenameFileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		utils.ReturnErrorResponse(c, err, fnRenameFile, utils.ErrSomethingWentWrong, "")
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if fileID == "" {
		utils.ReturnErrorResponse(c, utils.ErrFileIDRequired, fnRenameFile, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	var payload models.RenameFileRequest
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ReturnErrorResponse(c, utils.ErrInvalidJSON, fnRenameFile, utils.ErrSomethingWentWrong, "")
		return
	}

	newName := strings.TrimSpace(payload.NewName)

	oldName, err := h.storage.GetUserFileName(ctx, userID, fileID)
	if err != nil {
		slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to get user file name", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(c, err, fnRenameFile, utils.ErrSomethingWentWrong, "")
		return
	}
	if oldName == "" {
		c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
		return
	}

	if err := h.storage.UpdateUserFile(ctx, userID, fileID, newName); err != nil {
		slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to update user file mapping", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(c, err, fnRenameFile, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "userID", userID)

	renameRequest := &MMSpb.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	}
	if _, err = h.MMSClient.RenameFile(ctx, renameRequest); err != nil {
		if rbErr := h.storage.UpdateUserFile(ctx, userID, fileID, oldName); rbErr != nil {
			slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to rollback user file mapping", slog.Any(config.ErrorKey, rbErr))
		}
		status, msg := handlers.MapGRPCError(err, utils.ErrFailedToRenameFile.Error())
		slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to rename file in MMS", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
}
