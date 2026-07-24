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

func (h *UserFilesHandler) DownloadFileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		utils.ReturnErrorResponse(c, err, fnDownloadFile, utils.ErrSomethingWentWrong, "")
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if fileID == "" {
		utils.ReturnErrorResponse(c, utils.ErrFileIDRequired, fnDownloadFile, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	resp, err := h.client.DownloadFileGrpcHandler(ctx, userID, fileID)
	if err != nil {
		status, msg := handlers.MapGRPCError(err, utils.ErrFailedToDownloadFile.Error())
		slog.Error(handlers.LogPrefix(fnDownloadFile)+"failed to download file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.Data(http.StatusOK, utils.MimeTypeToString(resp.GetMimeType()), resp.GetContent())
}
