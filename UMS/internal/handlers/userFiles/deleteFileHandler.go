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

var deleteFileSuccessMsg = "file deleted successfully"

func (h *UserFilesHandler) DeleteFileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnDeleteFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	var fileID = strings.TrimSpace(c.Param("fileID"))
	if err := modelUtils.Validate.Var(fileID, "required,isValueEmpty,uuid"); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.ErrFileIDRequired, FnDeleteFile, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx := c.Request.Context()

	status, msg := h.client.DeleteFileGrpcHandler(ctx, userID, fileID)
	if status != http.StatusOK {
		c.JSON(status, gin.H{config.ErrorKey: msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": deleteFileSuccessMsg})
}
