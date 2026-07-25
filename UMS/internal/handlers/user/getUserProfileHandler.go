package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
)

func (h *UserHandler) GetUserProfileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, FnGetUserProfile, storage.ErrFailedToGetUserByID, config.NullString)
		return
	}

	ctx := c.Request.Context()

	user, err := h.storage.GetUserByID(ctx, userID)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, FnGetUserProfile, middlewareUtils.ErrSomethingWentWrong, user)
		return
	}

	user.Password = config.NullString

	c.JSON(http.StatusOK, gin.H{"user": user})
}
