package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUserHandler creates a new user account.
func (h *UserHandler) RegisterUserHandler(ctx *gin.Context) {
	var payload models.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		handlerErrors.ReturnErrorResponse(ctx, handlerErrors.ErrInvalidJSON, FnRegisterUser, handlerErrors.ErrFailedToRegisterUser)
		return
	}

	payload.Email = modelUtils.NormalizeEmail(payload.Email)

	if err := modelUtils.Validate.Struct(payload); err != nil {
		handlerErrors.ReturnErrorResponse(ctx, modelUtils.FieldErrors(err), FnRegisterUser, handlerErrors.ErrFailedToRegisterUser)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnRegisterUser)+"failed to generate password hash", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(ctx, err, FnRegisterUser, handlerErrors.ErrFailedToRegisterUser)
		return
	}

	password := string(hashedPassword)
	newUser, err := h.storage.CreateUser(ctx, &models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: password,
		Phone:    payload.Phone,
	})
	if err != nil {
		handlerErrors.ReturnErrorResponse(ctx, err, FnRegisterUser, handlerErrors.ErrFailedToRegisterUser)
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": newUser})
}
