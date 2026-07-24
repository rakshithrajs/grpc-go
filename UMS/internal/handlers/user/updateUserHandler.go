package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/handlers"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	userUpdatedMessage = "User profile updated successfully"
)

func (h *UserHandler) UpdateUserHandler(c *gin.Context) {
	id, err := handlers.GetUserIDFromGin(c)
	if err != nil {
		utils.ReturnErrorResponse(c, err, FnUpdateUserProfile, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx := c.Request.Context()

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ReturnErrorResponse(c, utils.ErrInvalidJSON, FnUpdateUserProfile, utils.ErrSomethingWentWrong, "")
		return
	}

	if req.Name == config.NullString && req.Email == config.NullString && req.Phone == config.NullString && req.Password == config.NullString {
		utils.ReturnErrorResponse(c, utils.ErrNoFieldsToUpdate, FnUpdateUserProfile, utils.ErrSomethingWentWrong, "")
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		utils.ReturnErrorResponse(c, utils.FieldErrors(err), FnUpdateUserProfile, utils.ErrSomethingWentWrong, "")
		return
	}

	if req.Password != config.NullString {
		user, err := h.storage.GetUserByID(ctx, id)
		if err != nil {
			utils.ReturnErrorResponse(c, err, FnUpdateUserProfile, utils.ErrSomethingWentWrong, user)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
			utils.ReturnErrorResponse(c, utils.ErrPasswordSameAsOldPassword, FnUpdateUserProfile, utils.ErrSomethingWentWrong, "")
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error(handlers.LogPrefix(FnUpdateUserProfile)+"failed to hash password", slog.Any(config.ErrorKey, err))
			utils.ReturnErrorResponse(c, err, FnUpdateUserProfile, utils.ErrSomethingWentWrong, "")
			return
		}
		hashed := string(hashedPassword)
		req.Password = hashed
	}

	if err := h.storage.UpdateUser(ctx, id, req); err != nil {
		utils.ReturnErrorResponse(c, err, FnUpdateUserProfile, utils.ErrSomethingWentWrong, "")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": userUpdatedMessage})
}
