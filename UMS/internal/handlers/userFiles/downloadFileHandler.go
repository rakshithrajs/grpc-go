package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

// DownloadFileHandler downloads a user file via the MMS client.
func (h *UserFilesHandler) DownloadFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnDownloadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	var uri models.FileIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrInvalidURI, FnDownloadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}
	if err := modelUtils.Validate.Struct(&uri); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnDownloadFile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	ctx := c.Request.Context()

	resp, err := h.client.DownloadFileGrpcClient(ctx, userID, uri.FileID)
	if err != nil {
		status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrFailedToDownloadFile.Error())
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.Data(http.StatusOK, handlerUtils.MimeTypeToString(resp.GetMimeType()), resp.GetContent())
}
