package handlers

import (
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"

	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
)

func (h *UserFilesHandler) UploadFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, fnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	fileHeader, err := c.FormFile(multipartFileField)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrFileIsRequired, fnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	if strings.TrimSpace(fileHeader.Filename) == config.NullString {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrFileNameRequired, fnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	openedFile, err := fileHeader.Open()
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(fnUploadFile)+"failed to open uploaded file", slog.Any(config.ErrorKey, err))
		handlerUtils.ReturnErrorResponse(c, err, fnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}
	defer openedFile.Close()

	content, err := io.ReadAll(openedFile)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(fnUploadFile)+"failed to read uploaded file", slog.Any(config.ErrorKey, err))
		handlerUtils.ReturnErrorResponse(c, err, fnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	if len(content) == 0 {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrEmptyFileContent, fnUploadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	file, status, msg := h.client.UploadFileGrpcHandler(ctx, userID, fileHeader.Filename, content)
	if status != http.StatusCreated {
		slog.Error(handlerUtils.LogPrefix(fnUploadFile)+"failed to upload file", slog.Any(config.ErrorKey, err))
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"file": file})
}
