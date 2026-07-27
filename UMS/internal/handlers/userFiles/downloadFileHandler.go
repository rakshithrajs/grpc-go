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

	fileID := strings.TrimSpace(c.Param("fileID"))
	if err := modelUtils.Validate.Struct(&modelUtils.FileIDPayload{FileID: fileID}); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnDownloadFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.client.DownloadFileGrpcHandler(ctx, userID, fileID)
	if err != nil {
		status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrFailedToDownloadFile.Error())
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.Data(http.StatusOK, handlerUtils.MimeTypeToString(resp.GetMimeType()), resp.GetContent())
}
