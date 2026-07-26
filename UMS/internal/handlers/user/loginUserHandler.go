package handlers

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"

	"golang.org/x/crypto/bcrypt"
)

func (h *UserHandler) LoginUserHandler(ctx *gin.Context) {

	var req models.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlerErrors.ReturnErrorResponse(ctx, handlerErrors.ErrInvalidJSON, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	req.Email = modelUtils.NormalizeEmail(req.Email)

	if err := modelUtils.Validate.Struct(&req); err != nil {
		fieldErrs := modelUtils.FieldErrors(err)
		handlerErrors.ReturnErrorResponse(ctx, fieldErrs, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	user, err := h.storage.GetUserByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		handlerErrors.ReturnErrorResponse(ctx, err, FnLoginUser, handlerErrors.ErrFailedToLoginUser, config.NullString)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		handlerErrors.ReturnErrorResponse(ctx, handlerErrors.ErrInvalidCredentials, FnLoginUser, handlerErrors.ErrFailedToLoginUser, config.NullString)
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnLoginUser)+"failed to get config", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(ctx, err, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	token, err := config.GenerateJWT(*user, cfg.JWTSecret)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnLoginUser)+"failed to generate JWT", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(ctx, err, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, config.NullString)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
