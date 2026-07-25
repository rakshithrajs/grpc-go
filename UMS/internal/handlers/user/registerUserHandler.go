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

func (h *UserHandler) RegisterUserHandler(ctx *gin.Context) {
	var payload models.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		handlerUtils.ReturnErrorResponse(ctx, handlerUtils.ErrInvalidJSON, FnRegisterUser, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	payload.Email = modelUtils.NormalizeEmail(payload.Email)

	if err := modelUtils.Validate.Struct(payload); err != nil {
		handlerUtils.ReturnErrorResponse(ctx, modelUtils.FieldErrors(err), FnRegisterUser, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnRegisterUser)+"failed to generate password hash", slog.Any(config.ErrorKey, err))
		handlerUtils.ReturnErrorResponse(ctx, err, FnRegisterUser, middlewareUtils.ErrSomethingWentWrong, "")
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
		handlerUtils.ReturnErrorResponse(ctx, err, FnRegisterUser, handlerUtils.ErrFailedToRegisterUser, "")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": newUser})
}
