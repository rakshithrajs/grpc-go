package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

const (
	AuthMiddlewareSource = "AuthMiddleware"
	LoginUserSource      = "LoginUserHandler"
	RegisterUserSource   = "RegisterUserHandler"
	UpdateUserSource     = "UpdateUserProfileHandler"
	GetUserProfileSource = "GetUserProfileHandler"
	UploadFileSource     = "UploadFile"
	DownloadFileSource   = "DownloadFile"
	ListFilesSource      = "ListFiles"
	RenameFileSource     = "RenameFile"
	DeleteFileSource     = "DeleteFile"
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
	ErrPasswordSameAsOldPassword = errors.New("new password is same as old password")
	ErrUserEmailAlreadyExists    = errors.New("user email already exists")
	ErrPhoneNumberAlreadyExists  = errors.New("phone number already exists")
	ErrFileIsRequired            = errors.New("file is required")
	ErrFileNameRequired          = errors.New("file name is required")
	ErrEmptyFileContent          = errors.New("file content is empty")
	ErrFailedToUploadFile        = errors.New("failed to upload file")
	ErrFailedToRenameFile        = errors.New("failed to rename file")
	ErrFailedToListFiles         = errors.New("failed to list files")
	ErrFailedToDownloadFile      = errors.New("failed to download file")
	ErrFailedToDeleteFile        = errors.New("failed to delete file")
	ErrUserFileAlreadyExists     = errors.New("user file mapping already exists")
	ErrFailedToRollback          = errors.New("failed to rollback changes")
	ErrFailedToCreateUser        = errors.New("failed to create user")
	ErrFailedToGetUserByID       = errors.New("failed to get user by ID")
	ErrFailedToGetUserByEmail    = errors.New("failed to get user by email")
	ErrFailedToUpdateUser        = errors.New("failed to update user")
	ErrFailedToCreateUserFile    = errors.New("failed to create user file mapping")
	ErrFailedToDeleteUserFile    = errors.New("failed to delete user file mapping")
	ErrFailedToUpdateUserFile    = errors.New("failed to update user file mapping")
	ErrFailedToListUserFiles     = errors.New("failed to list user files")
)

var errorResponseSpec = map[string]map[error]int{
	AuthMiddlewareSource: {
		ErrMissingAuthHeader:                  http.StatusUnauthorized,
		middlewareUtils.ErrMissingBearerToken: http.StatusUnauthorized,
		middlewareUtils.ErrInvalidToken:       http.StatusUnauthorized,
		middlewareUtils.ErrTokenExpired:       http.StatusUnauthorized,
	},
	LoginUserSource: {
		ErrInvalidJSON:        http.StatusBadRequest,
		ErrInvalidCredentials: http.StatusUnauthorized,
		ErrFailedToLoginUser:  http.StatusInternalServerError,
	},
	RegisterUserSource: {
		ErrInvalidJSON:              http.StatusBadRequest,
		ErrUserEmailAlreadyExists:   http.StatusConflict,
		ErrPhoneNumberAlreadyExists: http.StatusConflict,
		ErrFailedToRegisterUser:     http.StatusInternalServerError,
	},
	UpdateUserSource: {
		ErrInvalidJSON:               http.StatusBadRequest,
		ErrNoFieldsToUpdate:          http.StatusBadRequest,
		ErrPasswordSameAsOldPassword: http.StatusBadRequest,
	},
	GetUserProfileSource: {
		ErrUserNotFound: http.StatusOK,
	},
	UploadFileSource: {
		ErrFileIsRequired:   http.StatusBadRequest,
		ErrFileNameRequired: http.StatusBadRequest,
		ErrEmptyFileContent: http.StatusBadRequest,
	},
	DownloadFileSource: {
		modelUtils.ErrFileIDRequired: http.StatusBadRequest,
	},
	ListFilesSource: {
		ErrFailedToListFiles: http.StatusInternalServerError,
	},
	RenameFileSource: {
		ErrInvalidJSON:               http.StatusBadRequest,
		modelUtils.ErrFileIDRequired: http.StatusBadRequest,
	},
	DeleteFileSource: {
		modelUtils.ErrFileIDRequired: http.StatusBadRequest,
	},
}

var commonInternalServerErrors = []error{
	ErrFailedToUploadFile,
	ErrFailedToDownloadFile,
	ErrFailedToRenameFile,
	ErrFailedToDeleteFile,
	ErrFailedToListFiles,
	ErrFailedToRegisterUser,
	ErrFailedToLoginUser,
	ErrFailedToRollback,
	ErrFailedToCreateUser,
	ErrFailedToGetUserByID,
	ErrFailedToGetUserByEmail,
	ErrFailedToUpdateUser,
	ErrFailedToCreateUserFile,
	ErrFailedToDeleteUserFile,
	ErrFailedToUpdateUserFile,
	ErrFailedToListUserFiles,
}

func ReturnErrorResponse(c *gin.Context, err any, source string, defaultMsg error, data any) {
	if e, ok := err.(error); ok {
		if isAuthError(e) {
			c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: e.Error()})
			return
		}

		if isCommonInternalError(e) {
			c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: e.Error()})
			return
		}

		if spec, ok := errorResponseSpec[source]; ok {
			if status, ok := spec[e]; ok {
				if status == http.StatusOK {
					c.JSON(http.StatusOK, gin.H{"user": data})
					return
				}
				c.JSON(status, gin.H{config.ErrorKey: e.Error()})
				return
			}
		}

		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: defaultMsg.Error()})
		return
	}

	errs, _ := err.(map[string]string)
	returnMultipleErrorResponse(c, errs, source, defaultMsg)
}

func isAuthError(e error) bool {
	return errors.Is(e, ErrMissingAuthHeader) ||
		errors.Is(e, middlewareUtils.ErrMissingBearerToken) ||
		errors.Is(e, middlewareUtils.ErrInvalidToken) ||
		errors.Is(e, middlewareUtils.ErrTokenExpired) ||
		errors.Is(e, ErrUnauthorized)
}

func isCommonInternalError(e error) bool {
	for _, candidate := range commonInternalServerErrors {
		if errors.Is(e, candidate) {
			return true
		}
	}
	return false
}

func returnMultipleErrorResponse(c *gin.Context, errs map[string]string, source string, _ error) {
	isLogin := source == LoginUserSource
	isRequiredField := errs["email"] == modelUtils.ErrEmailRequired.Error() || errs["password"] == modelUtils.ErrPasswordRequired.Error()

	if !isLogin || (isLogin && isRequiredField) {
		c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: errs})
		return
	}

	c.JSON(http.StatusUnauthorized, gin.H{config.ErrorKey: ErrInvalidCredentials.Error()})
}
