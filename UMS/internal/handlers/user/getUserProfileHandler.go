package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func (h *UserHandler) GetUserProfileHandler(c *gin.Context) {
	userID, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		utils.ReturnErrorResponse(c, err, FnGetUserProfile, storage.ErrFailedToGetUserByID, "")
		return
	}

	ctx := c.Request.Context()

	user, err := h.storage.GetUserByID(ctx, userID)
	if err != nil {
		utils.ReturnErrorResponse(c, err, FnGetUserProfile, utils.ErrSomethingWentWrong, user)
		return
	}

	user.Password = config.NullString

	c.JSON(http.StatusOK, gin.H{"user": user})
}
