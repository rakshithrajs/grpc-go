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

// LoginUserHandler authenticates a user and returns a JWT access token.
func (h *UserHandler) LoginUserHandler(ctx *gin.Context) {

	var req models.LoginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		handlerErrors.ReturnErrorResponse(ctx, handlerErrors.ErrInvalidJSON, FnLoginUser, handlerErrors.ErrSomethingWentWrong)
		return
	}

	req.Email = modelUtils.NormalizeEmail(req.Email)

	if err := modelUtils.Validate.Struct(&req); err != nil {
		fieldErrs := modelUtils.FieldErrors(err)
		handlerErrors.ReturnErrorResponse(ctx, fieldErrs, FnLoginUser, handlerErrors.ErrSomethingWentWrong)
		return
	}

	user, err := h.storage.GetUserByEmail(ctx.Request.Context(), req.Email)
	if err != nil {
		handlerErrors.ReturnErrorResponse(ctx, err, FnLoginUser, handlerErrors.ErrFailedToLoginUser)
		return
	}
	if user == nil {
		handlerErrors.ReturnErrorResponse(ctx, handlerErrors.ErrInvalidCredentials, FnLoginUser, handlerErrors.ErrFailedToLoginUser)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		handlerErrors.ReturnErrorResponse(ctx, handlerErrors.ErrInvalidCredentials, FnLoginUser, handlerErrors.ErrFailedToLoginUser)
		return
	}

	token, err := h.tmsClient.GenerateToken(ctx.Request.Context(), user.ID)
	if err != nil {
		slog.Error(handlerUtils.LogPrefix(FnLoginUser)+"failed to generate token", slog.Any(config.ErrorKey, err))
		handlerErrors.ReturnErrorResponse(ctx, err, FnLoginUser, handlerErrors.ErrSomethingWentWrong)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": token})
}
