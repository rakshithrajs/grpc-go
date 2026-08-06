package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
)

// GetUserProfileHandler returns the authenticated user's profile.
func (h *UserHandler) GetUserProfileHandler(c *gin.Context) {
	userID, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnGetUserProfile, storage.ErrFailedToGetUserByID)
		return
	}

	ctx := c.Request.Context()

	user, err := h.storage.GetUserByID(ctx, userID)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnGetUserProfile, handlerErrors.ErrSomethingWentWrong)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
