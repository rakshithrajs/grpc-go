package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

func (h *UserFilesHandler) DownloadFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnDownloadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	var fileID = strings.TrimSpace(c.Param("fileID"))
	if err := modelUtils.Validate.Var(fileID, "required,isValueEmpty,uuid"); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.ErrFileIDRequired, FnDownloadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	resp, status, msg := h.client.DownloadFileGrpcHandler(ctx, userID, fileID)
	if status != http.StatusOK {
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.Data(http.StatusOK, handlerUtils.MimeTypeToString(resp.GetMimeType()), resp.GetContent())
}
