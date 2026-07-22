package handlers

import (
	"io"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"google.golang.org/grpc/metadata"
)

func (h *UserFilesHandler) UploadFileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: err.Error()})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: handlers.ErrFileIsRequired.Error()})
		return
	}

	openedFile, err := fileHeader.Open()
	if err != nil {
		slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to open uploaded file", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}
	defer openedFile.Close()

	content, err := io.ReadAll(openedFile)
	if err != nil {
		slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to read uploaded file", slog.Any(config.ErrorKey, err))
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}

	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "x-user-id", userID)
	resp, err := h.MMSClient.UploadFile(ctx, &MMSpb.UploadFileRequest{
		FileName: fileHeader.Filename,
		Content:  content,
	})
	if err != nil {
		status, msg := handlers.MapGRPCError(err, handlers.ErrFailedToUploadFile.Error())
		slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to upload file to MMS", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	if resp.GetFile() == nil || resp.GetFile().GetID() == "" {
		slog.Error(handlers.LogPrefix(fnUploadFile) + "MMS returned empty file response")
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrSomethingWentWrong.Error()})
		return
	}

	fileID := resp.GetFile().GetID()
	if err := h.storage.CreateUserFile(c.Request.Context(), userID, fileID, resp.File.FileName); err != nil {
		slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to save user file mapping", slog.Any(config.ErrorKey, err))
		if _, delErr := h.MMSClient.DeleteFile(ctx, &MMSpb.DeleteFileRequest{FileID: fileID}); delErr != nil {
			slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to compensate MMS upload", slog.Any(config.ErrorKey, delErr))
		}
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: handlers.ErrFailedToUploadFile.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"file": resp.GetFile()})
}
