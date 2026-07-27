package handlers

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

var deleteFileSuccessMsg = "file deleted successfully"

func (h *UserFilesHandler) DeleteFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnDeleteFile, handlerErrors.ErrFailedToDeleteFile, config.NullString)
		return
	}

	fileID := strings.TrimSpace(c.Param("fileID"))
	if err := modelUtils.Validate.Struct(&modelUtils.FileIDPayload{FileID: fileID}); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnDeleteFile, handlerErrors.ErrFailedToDeleteFile, config.NullString)
		return
	}
	ctx := c.Request.Context()

	if err = h.client.DeleteFileGrpcClient(ctx, userID, fileID); err != nil {
		status, errMsg := handlerUtils.MapGRPCError(err, handlerErrors.ErrFailedToDeleteFile.Error())
		c.JSON(status, gin.H{"error": errMsg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": deleteFileSuccessMsg})
}
