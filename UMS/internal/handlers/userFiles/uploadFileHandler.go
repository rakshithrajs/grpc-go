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
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
)

func (h *UserFilesHandler) UploadFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	fileHeader, err := c.FormFile(multipartFileField)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrFileIsRequired, FnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	if strings.TrimSpace(fileHeader.Filename) == config.NullString {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrFileNameRequired, FnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	openedFile, err := fileHeader.Open()
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnUploadFile)+"failed to open uploaded file", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(c, err, FnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}
	defer openedFile.Close()

	content, err := io.ReadAll(openedFile)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnUploadFile)+"failed to read uploaded file", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(c, err, FnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	if len(content) == 0 {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrEmptyFileContent, FnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
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
