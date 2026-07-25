package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func (h *UserFilesHandler) UploadFileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		utils.ReturnErrorResponse(c, err, fnUploadFile, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	fileHeader, err := c.FormFile(multipartFileField)
	if err != nil {
		utils.ReturnErrorResponse(c, utils.ErrFileIsRequired, fnUploadFile, utils.ErrSomethingWentWrong, "")
		return
	}

	if strings.TrimSpace(fileHeader.Filename) == config.NullString {
		utils.ReturnErrorResponse(c, utils.ErrFileNameRequired, fnUploadFile, utils.ErrSomethingWentWrong, "")
		return
	}

	openedFile, err := fileHeader.Open()
	if err != nil {
		slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to open uploaded file", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(c, err, fnUploadFile, utils.ErrSomethingWentWrong, "")
		return
	}
	defer openedFile.Close()

	content, err := io.ReadAll(openedFile)
	if err != nil {
		slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to read uploaded file", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(c, err, fnUploadFile, utils.ErrSomethingWentWrong, "")
		return
	}

	if len(content) == 0 {
		utils.ReturnErrorResponse(c, utils.ErrEmptyFileContent, fnUploadFile, utils.ErrSomethingWentWrong, "")
		return
	}

	file, err := h.client.UploadFileGrpcHandler(ctx, userID, fileHeader.Filename, content)
	if err != nil {
		status, msg := handlers.MapGRPCError(err, utils.ErrFailedToUploadFile.Error())
		slog.Error(handlers.LogPrefix(fnUploadFile)+"failed to upload file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"file": file})
}
