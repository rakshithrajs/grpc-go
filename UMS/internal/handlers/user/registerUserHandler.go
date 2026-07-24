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

func (a *UMSHandler) RegisterUserHandler(ctx *gin.Context) {
	var payload models.RegisterUserRequest
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		utils.ReturnErrorResponse(ctx, utils.ErrInvalidJSON, FnRegisterUser, utils.ErrSomethingWentWrong, "")
		return
	}

	payload.Email = utils.NormalizeEmail(payload.Email)

	if err := utils.Validate.Struct(payload); err != nil {
		utils.ReturnErrorResponse(ctx, utils.FieldErrors(err), FnRegisterUser, utils.ErrSomethingWentWrong, "")
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(handlers.LogPrefix(FnRegisterUser)+"failed to generate password hash", slog.Any(config.ErrorKey, err))
		utils.ReturnErrorResponse(ctx, err, FnRegisterUser, utils.ErrSomethingWentWrong, "")
		return
	}

	password := string(hashedPassword)
	newUser, err := a.storage.CreateUser(ctx, &models.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: password,
		Phone:    payload.Phone,
	})
	if err != nil {
		utils.ReturnErrorResponse(ctx, err, FnRegisterUser, utils.ErrFailedToRegisterUser, "")
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"user": newUser})
}
