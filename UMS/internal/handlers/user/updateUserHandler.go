package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	userUpdatedMessage = "User profile updated successfully"
)

// UpdateUserHandler updates the authenticated user's profile.
func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	id, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
		return
	}

	ctx := c.Request.Context()

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrInvalidJSON, FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
		return
	}

	if req.Name == config.NullString && req.Email == config.NullString && req.Phone == config.NullString && req.Password == config.NullString {
		handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrNoFieldsToUpdate, FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
		return
	}

	if err := modelUtils.Validate.Struct(req); err != nil {
		handlerErrors.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
		return
	}

	req.Email = modelUtils.NormalizeEmail(req.Email)

	if req.Password != config.NullString {
		user, err := h.storage.GetUserByID(ctx, id)
		if err != nil {
			handlerErrors.ReturnErrorResponse(c, err, FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
			return
		}
		if user == nil {
			c.JSON(http.StatusOK, gin.H{"message": userUpdatedMessage})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
			handlerErrors.ReturnErrorResponse(c, handlerErrors.ErrPasswordSameAsOldPassword, FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error(handlerUtils.LogPrefix(FnUpdateUserProfile)+"failed to hash password", slog.Any(config.ErrorKey, err))
			handlerErrors.ReturnErrorResponse(c, err, FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
			return
		}
		hashed := string(hashedPassword)
		req.Password = hashed
	}

	if err := h.storage.UpdateUser(ctx, id, req); err != nil {
		handlerErrors.ReturnErrorResponse(c, err, FnUpdateUserProfile, storage.ErrFailedToUpdateUser)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": userUpdatedMessage})
}
