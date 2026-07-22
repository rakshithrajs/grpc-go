package handlers

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var (
	userUpdatedMessage = "User profile updated successfully"
)

func (a *UMSHandler) UpdateUserHandler(c *gin.Context) {
	_, err := GetUser(c, "UpdateUserHandler")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: err.Error()})
		return
	}

	id := strings.TrimSpace(c.Param("id"))
	if err := utils.Validate.Var(id, "required,uuid"); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: ErrInvalidID.Error()})
		return
	}

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: ErrInvalidJSON.Error()})
		return
	}

	if req.Password == config.NullString && req.Phone == config.NullString {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: ErrNoFieldsToUpdate.Error()})
		return
	}

	if err := utils.Validate.Struct(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: utils.FieldErrors(err)})
		return
	}

	if req.Name == config.NullString && req.Email == config.NullString && req.Phone == config.NullString && req.Password == config.NullString {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: ErrNoFieldsToUpdate.Error()})
		return
	}

	if req.Password != config.NullString {
		user, err := a.storage.GetUserByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				c.JSON(http.StatusOK, gin.H{"message": userUpdatedMessage})
				return
			}
			slog.Error(logPrefix(fnUpdateUserProfile)+"failed to get user by id", slog.Any(config.ErrorKey, err))
			c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: ErrSomethingWentWrong.Error()})
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err == nil {
			c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: storage.ErrPasswordSameAsOldPassword.Error()})
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error(logPrefix(fnUpdateUserProfile)+"failed to hash password", slog.Any(config.ErrorKey, err))
			c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: ErrSomethingWentWrong.Error()})
			return
		}
		hashed := string(hashedPassword)
		req.Password = hashed
	}

	if err := a.storage.UpdateUser(c.Request.Context(), id, req); err != nil {
		if HandleDomainError(c, err) {
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: ErrSomethingWentWrong.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": userUpdatedMessage})
}
