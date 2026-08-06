package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
)

// UploadFileHandler uploads a user file via the MMS client.
func (h *UserFilesHandler) UploadFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnUploadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	ctx := c.Request.Context()

	fileHeader, err := c.FormFile(multipartFileField)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrFileIsRequired, FnUploadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	if strings.TrimSpace(fileHeader.Filename) == config.NullString {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrFileNameRequired, FnUploadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	openedFile, err := fileHeader.Open()
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnUploadFile)+"failed to open uploaded file", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(c, err, FnUploadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}
	defer func() {
		if err := openedFile.Close(); err != nil {
			slog.Error(handlerUtils.LogPrefix(FnUploadFile)+"failed to close uploaded file", slog.Any(config.ErrorKey, err))
		}
	}()

	content, err := io.ReadAll(openedFile)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnUploadFile)+"failed to read uploaded file", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(c, err, FnUploadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	if len(content) == 0 {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrEmptyFileContent, FnUploadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	file, err := h.client.UploadFileGrpcClient(ctx, userID, fileHeader.Filename, content)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnUploadFile)+"failed to upload file", slog.Any(config.ErrorKey, err))
		status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrFailedToUploadFile.Error())
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"file": file})
}
