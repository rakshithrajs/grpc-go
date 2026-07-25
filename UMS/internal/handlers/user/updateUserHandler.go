package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"

	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	userUpdatedMessage = "User profile updated successfully"
)

func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	id, err := handlerUtils.GetUserIDFromGin(c)
	if err != nil {
		handlerUtils.ReturnErrorResponse(c, err, FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrInvalidJSON, FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	if req.Name == config.NullString && req.Email == config.NullString && req.Phone == config.NullString && req.Password == config.NullString {
		handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrNoFieldsToUpdate, FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	if err := modelUtils.Validate.Struct(req); err != nil {
		handlerUtils.ReturnErrorResponse(c, modelUtils.FieldErrors(err), FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	req.Email = modelUtils.NormalizeEmail(req.Email)

	if req.Password != config.NullString {
		user, err := h.storage.GetUserByID(ctx, id)
		if err != nil {
			handlerUtils.ReturnErrorResponse(c, err, FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, user)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
			handlerUtils.ReturnErrorResponse(c, handlerUtils.ErrPasswordSameAsOldPassword, FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, "")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error(handlerUtils.LogPrefix(FnUpdateUserProfile)+"failed to hash password", slog.Any(config.ErrorKey, err))
			handlerUtils.ReturnErrorResponse(c, err, FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, "")
			return
		}
		hashed := string(hashedPassword)
		req.Password = hashed
	}

	if err := h.storage.UpdateUser(ctx, id, req); err != nil {
		handlerUtils.ReturnErrorResponse(c, err, FnUpdateUserProfile, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": userUpdatedMessage})
}
