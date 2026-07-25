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

func (h *UserFilesHandler) DownloadFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, fnDownloadFile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if fileID == "" {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrFileIDRequired, fnDownloadFile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	resp, err := h.client.DownloadFileGrpcHandler(ctx, userID, fileID)
	if err != nil {
		status, msg := grpcUtils.MapGRPCError(err, handlerUtils.ErrFailedToDownloadFile.Error())
		slog.Error(handlerUtils.LogPrefix(fnDownloadFile)+"failed to download file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.Data(http.StatusOK, handlerUtils.MimeTypeToString(resp.GetMimeType()), resp.GetContent())
}
