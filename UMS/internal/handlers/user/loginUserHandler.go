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

func (h *UserHandler) LoginUserHandler(ctx *gin.Context) {

	var req models.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlerUtils.ReturnErrorResponse(ctx, handlerUtils.ErrInvalidJSON, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	req.Email = modelUtils.NormalizeEmail(req.Email)

	if err := modelUtils.Validate.Struct(&req); err != nil {
		fieldErrs := modelUtils.FieldErrors(err)
		handlerUtils.ReturnErrorResponse(ctx, fieldErrs, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	user, err := h.storage.GetUserByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		handlerUtils.ReturnErrorResponse(ctx, err, FnLoginUser, handlerUtils.ErrFailedToLoginUser, "")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		handlerUtils.ReturnErrorResponse(ctx, handlerUtils.ErrInvalidCredentials, FnLoginUser, handlerUtils.ErrFailedToLoginUser, "")
		return
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnLoginUser)+"failed to get config", slog.Any(config.ErrorKey, err))
		handlerUtils.ReturnErrorResponse(ctx, err, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	token, err := config.GenerateJWT(*user, cfg.JWTSecret)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnLoginUser)+"failed to generate JWT", slog.Any(config.ErrorKey, err))
		handlerUtils.ReturnErrorResponse(ctx, err, FnLoginUser, middlewareUtils.ErrSomethingWentWrong, "")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
