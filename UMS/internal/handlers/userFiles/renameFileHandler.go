package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
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
		utils.ReturnErrorResponse(c, utils.FieldErrors(err), fnRenameFile, utils.ErrSomethingWentWrong, "")
		return
	}

	newName := strings.TrimSpace(payload.NewName)

	if err := h.client.RenameFileGrpcHandler(ctx, userID, fileID, newName); err != nil {
		status, msg := handlers.MapGRPCError(err, utils.ErrFailedToRenameFile.Error())
		slog.Error(handlers.LogPrefix(fnRenameFile)+"failed to rename file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fileRenamedSuccessfully})
}
