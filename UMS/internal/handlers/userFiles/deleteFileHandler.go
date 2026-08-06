package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

var deleteFileSuccessMsg = "file deleted successfully"

// DeleteFileHandler deletes a user file via the MMS client.
func (h *UserFilesHandler) DeleteFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnDeleteFile, handlerErrors.ErrFailedToDeleteFile)
		return
	}

	var uri models.FileIDURI
	if err := c.ShouldBindUri(&uri); err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrInvalidURI, FnDeleteFile, handlerErrors.ErrFailedToDeleteFile)
		return
	}
	if err := modelUtils.Validate.Struct(&uri); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnDeleteFile, handlerErrors.ErrFailedToDeleteFile)
		return
	}

	ctx := c.Request.Context()

	if err = h.client.DeleteFileGrpcClient(ctx, userID, uri.FileID); err != nil {
		status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrFailedToDeleteFile.Error())
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": deleteFileSuccessMsg})
}
