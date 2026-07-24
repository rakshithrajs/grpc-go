package utils

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
)

var (
	ErrMissingAuthHeader         = errors.New("missing Authorization header")
	ErrUnauthorized              = errors.New("unauthorized")
	ErrUserNotFound              = errors.New("user not found")
	ErrEmailNotFound             = errors.New("email not found")
	ErrInvalidJSON               = errors.New("invalid JSON")
	ErrFailedToRegisterUser      = errors.New("failed to register user")
	ErrInvalidCredentials        = errors.New("invalid credentials")
	ErrFailedToLoginUser         = errors.New("failed to login user")
	ErrNoFieldsToUpdate          = errors.New("no fields to update")
	ErrInvalidID                 = errors.New("invalid ID")
	ErrPasswordSameAsOldPassword = errors.New("new password is same as old password")
	ErrUserEmailAlreadyExists    = errors.New("user email already exists")
	ErrPhoneNumberAlreadyExists  = errors.New("phone number already exists")
	ErrFileIDRequired            = errors.New("file ID is required")
	ErrFileIsRequired            = errors.New("file is required")
	ErrFailedToUploadFile        = errors.New("failed to upload file")
	ErrFailedToRenameFile        = errors.New("failed to rename file")
	ErrFailedToListFiles         = errors.New("failed to list files")
	ErrFailedToDownloadFile      = errors.New("failed to download file")
	ErrFailedToDeleteFile        = errors.New("failed to delete file")
	ErrFileRenamed               = errors.New("file renamed successfully")
	ErrFileDeleted               = errors.New("file deleted successfully")
	ErrFileNotFound              = errors.New("file not found")
	ErrUserFileAlreadyExists     = errors.New("user file mapping already exists")
	ErrFailedToCreateUserFile    = errors.New("failed to create user file mapping")
	ErrFailedToDeleteUserFile    = errors.New("failed to delete user file mapping")
	ErrFailedToUpdateUserFile    = errors.New("failed to update user file mapping")
	ErrFailedToListUserFiles     = errors.New("failed to list user files")
	ErrFailedToCreateUser        = errors.New("failed to create user")
	ErrFailedToGetUserByID       = errors.New("failed to get user by ID")
	ErrFailedToGetUserByEmail    = errors.New("failed to get user by email")
	ErrFailedToUpdateUser        = errors.New("failed to update user")
	ErrFailedToRollback          = errors.New("failed to rollback changes")
)

func ReturnErrorResponse(c *gin.Context, err any, source string, defaultMsg error, data any) {
	if e, ok := err.(error); ok {
		switch {
		case errors.Is(e, ErrMissingAuthHeader), errors.Is(e, ErrMissingBearerToken), errors.Is(e, ErrInvalidToken), errors.Is(e, ErrTokenExpired), errors.Is(e, ErrUnauthorized):
			c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: e.Error()})
			return
		case errors.Is(e, ErrInvalidJSON), errors.Is(e, ErrNoFieldsToUpdate), errors.Is(e, ErrPasswordSameAsOldPassword), errors.Is(e, ErrFileIDRequired), errors.Is(e, ErrFileIsRequired):
			c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: e.Error()})
			return
		case errors.Is(e, ErrInvalidCredentials), errors.Is(e, ErrEmailNotFound):
			c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: ErrInvalidCredentials.Error()})
			return
		case errors.Is(e, ErrUserEmailAlreadyExists), errors.Is(e, ErrPhoneNumberAlreadyExists), errors.Is(e, ErrUserFileAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{config.ErrorKey: e.Error()})
			return
		case errors.Is(e, ErrUserNotFound):
			c.JSON(http.StatusOK, gin.H{"user": data})
			return
		case errors.Is(e, ErrFailedToRollback):
			c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: ErrFailedToRollback.Error()})
			return
		case errors.Is(e, ErrFailedToCreateUser), errors.Is(e, ErrFailedToGetUserByID), errors.Is(e, ErrFailedToGetUserByEmail),
			errors.Is(e, ErrFailedToUpdateUser), errors.Is(e, ErrFailedToCreateUserFile), errors.Is(e, ErrFailedToDeleteUserFile),
			errors.Is(e, ErrFailedToUpdateUserFile), errors.Is(e, ErrFailedToListUserFiles),
			errors.Is(e, ErrFailedToUploadFile), errors.Is(e, ErrFailedToDownloadFile), errors.Is(e, ErrFailedToRenameFile),
			errors.Is(e, ErrFailedToDeleteFile), errors.Is(e, ErrFailedToListFiles), errors.Is(e, ErrFailedToRegisterUser):
			c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: e.Error()})
			return
		default:
			c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: defaultMsg.Error()})
			return
		}
	}

	errs, _ := err.(map[string]string)
	returnMultipleErrorResponse(c, errs, source, defaultMsg)
}

func returnMultipleErrorResponse(c *gin.Context, errs map[string]string, source string, _ error) {
	isLogin := source == "LoginUserHandler"
	isRequiredField := errs["email"] == ErrEmailRequired.Error() || errs["password"] == ErrPasswordRequired.Error()

	if !isLogin || (isLogin && isRequiredField) {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: errs})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: ErrInvalidCredentials.Error()})
}
