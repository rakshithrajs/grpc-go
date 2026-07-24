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

func (a *UMSHandler) LoginUserHandler(ctx *gin.Context) {

	var req models.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.ReturnErrorResponse(ctx, utils.ErrInvalidJSON, FnLoginUser, utils.ErrSomethingWentWrong, "")
		return
	}

	req.Email = utils.NormalizeEmail(req.Email)

	if err := utils.Validate.Struct(&req); err != nil {
		fieldErrs := utils.FieldErrors(err)
		utils.ReturnErrorResponse(ctx, fieldErrs, FnLoginUser, utils.ErrSomethingWentWrong, "")
		return
	}

	user, err := a.storage.GetUserByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		utils.ReturnErrorResponse(ctx, err, FnLoginUser, utils.ErrFailedToLoginUser, "")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		utils.ReturnErrorResponse(ctx, utils.ErrInvalidCredentials, FnLoginUser, utils.ErrFailedToLoginUser, "")
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(handlers.LogPrefix(FnLoginUser)+"failed to get config", slog.Any("error", err))
		utils.ReturnErrorResponse(ctx, err, FnLoginUser, utils.ErrSomethingWentWrong, "")
		return
	}

	token, err := config.GenerateJWT(*user, cfg.JWTSecret)
	if err != nil {
		slog.Error(handlers.LogPrefix(FnLoginUser)+"failed to generate JWT", slog.Any("error", err))
		utils.ReturnErrorResponse(ctx, err, FnLoginUser, utils.ErrSomethingWentWrong, "")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
