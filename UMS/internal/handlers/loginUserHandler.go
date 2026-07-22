package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

func (a *UMSHandler) LoginUserHandler(ctx *gin.Context) {

	var req models.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: ErrInvalidJSON.Error()})
		return
	}

	req.Email = utils.NormalizeEmail(req.Email)

	if err := utils.Validate.Struct(&req); err != nil {
		fieldErrs := utils.FieldErrors(err)
		if fieldErrs["email"] == utils.ErrEmailRequired.Error() || fieldErrs["password"] == utils.ErrPasswordRequired.Error() {
			ctx.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: fieldErrs})
			return
		}
		ctx.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: ErrInvalidCredentials.Error()})
		return
	}

	user, err := a.storage.GetUserByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		if errors.Is(err, storage.ErrEmailNotFound) {
			ctx.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: ErrInvalidCredentials.Error()})
			return
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: ErrFailedToLoginUser.Error()})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: ErrInvalidCredentials.Error()})
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(logPrefix(fnLoginUser)+"failed to get config", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: ErrSomethingWentWrong.Error()})
		return
	}

	token, err := config.GenerateJWT(*user, cfg.JWTSecret)
	if err != nil {
		slog.Error(logPrefix(fnLoginUser)+"failed to generate JWT", slog.Any("error", err))
		ctx.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: ErrSomethingWentWrong.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{config.ErrorKey: nil, "token": token})
}
