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

func (h *UserFilesHandler) DownloadFileHandler(c *gin.Context) {
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

	ctx := metadata.AppendToOutgoingContext(c.Request.Context(), "x-user-id", userID)
	resp, err := h.MMSClient.DownloadFile(ctx, &MMSpb.DownloadFileRequest{FileID: fileID})
	if err != nil {
		status, msg := handlers.MapGRPCError(err, handlers.ErrFailedToDownloadFile.Error())
		slog.Error(handlers.LogPrefix(fnDownloadFile)+"failed to download file from MMS", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.Data(http.StatusOK, resp.GetMimeType(), resp.GetContent())
}
