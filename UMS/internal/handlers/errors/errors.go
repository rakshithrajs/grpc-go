package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

// source identifiers for error handling
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
	// message: missing Authorization header
	ErrMissingAuthHeader = errors.New("missing Authorization header")

	// message: unauthorized
	ErrUnauthorized = errors.New("unauthorized")

	// message: user not found
	ErrUserNotFound = errors.New("user not found")

	// message: email not found
	ErrEmailNotFound = errors.New("email not found")

	// message: invalid JSON
	ErrInvalidJSON = errors.New("invalid JSON")

	// message: failed to register user
	ErrFailedToRegisterUser = errors.New("failed to register user")

	// message: invalid credentials
	ErrInvalidCredentials = errors.New("invalid credentials")

	// message: failed to login user
	ErrFailedToLoginUser = errors.New("failed to login user")

	// message: no fields to update
	ErrNoFieldsToUpdate = errors.New("no fields to update")

	// message: new password is same as old password
	ErrPasswordSameAsOldPassword = errors.New("new password is same as old password")

	// message: user email already exists
	ErrUserEmailAlreadyExists = errors.New("user email already exists")

	// message: phone number already exists
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")

	// message: file is required
	ErrFileIsRequired = errors.New("file is required")

	// message: file name is required
	ErrFileNameRequired = errors.New("file name is required")

	// message: file content is empty
	ErrEmptyFileContent = errors.New("file content is empty")

	// message: failed to upload file
	ErrFailedToUploadFile = errors.New("failed to upload file")

	// message: failed to rename file
	ErrFailedToRenameFile = errors.New("failed to rename file")

	// message: failed to list file
	ErrFailedToListFiles = errors.New("failed to list files")

	// message: failed to download file
	ErrFailedToDownloadFile = errors.New("failed to download file")

	// message: failed to delete file
	ErrFailedToDeleteFile = errors.New("failed to delete file")

	// message: user file mapping already exists
	ErrUserFileAlreadyExists = errors.New("user file mapping already exists")

	// message: failed to rollback changes
	ErrFailedToRollback = errors.New("failed to rollback changes")

	// message: failed to create user
	ErrFailedToCreateUser = errors.New("failed to create user")

	// message: failed to get user by ID
	ErrFailedToGetUserByID = errors.New("failed to get user by ID")

	// message: failed to get user by email
	ErrFailedToGetUserByEmail = errors.New("failed to get user by email")

	// message: failed to update user
	ErrFailedToUpdateUser = errors.New("failed to update user")

	// message: failed to create user file mapping
	ErrFailedToCreateUserFile = errors.New("failed to create user file mapping")

	// message: failed to delete user file mapping
	ErrFailedToDeleteUserFile = errors.New("failed to delete user file mapping")

	// message: failed to update user file mapping
	ErrFailedToUpdateUserFile = errors.New("failed to update user file mapping")

	// message: failed to list user files
	ErrFailedToListUserFiles = errors.New("failed to list user files")
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
		ErrEmailNotFound:      http.StatusUnauthorized,
		ErrFailedToLoginUser:  http.StatusInternalServerError,
	},
	RegisterUserSource: {
		ErrInvalidJSON:              http.StatusBadRequest,
		ErrUserEmailAlreadyExists:   http.StatusConflict,
		ErrPhoneNumberAlreadyExists: http.StatusConflict,
		ErrFailedToRegisterUser:     http.StatusInternalServerError,
	},
	UpdateUserSource: {
		ErrUserNotFound:              http.StatusOK,
		ErrInvalidJSON:               http.StatusBadRequest,
		ErrNoFieldsToUpdate:          http.StatusBadRequest,
		ErrPasswordSameAsOldPassword: http.StatusBadRequest,
		ErrUserEmailAlreadyExists:    http.StatusConflict,
		ErrPhoneNumberAlreadyExists:  http.StatusConflict,
		ErrFailedToUpdateUser:        http.StatusInternalServerError,
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
		ErrInvalidJSON:                   http.StatusBadRequest,
		modelUtils.ErrFileIDRequired:     http.StatusBadRequest,
		modelUtils.ErrFileIDInvalidUUID:  http.StatusBadRequest,
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

// ReturnErrorResponse sends an error response based on the provided error and source. It checks for specific error types and returns appropriate HTTP status codes and messages.
func ReturnErrorResponse(c *gin.Context, err any, source string, defaultError error, data any) {
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

		c.JSON(http.StatusInternalServerError, gin.H{config.ErrorKey: defaultError.Error()})
		return
	}

	errs, _ := err.(map[string]string)
	c.JSON(http.StatusBadRequest, gin.H{config.ErrorKey: errs})
}

// isAuthError checks if the provided error is related to authentication issues.
func isAuthError(e error) bool {
	return errors.Is(e, ErrMissingAuthHeader) ||
		errors.Is(e, middlewareUtils.ErrMissingBearerToken) ||
		errors.Is(e, middlewareUtils.ErrInvalidToken) ||
		errors.Is(e, middlewareUtils.ErrTokenExpired) ||
		errors.Is(e, ErrUnauthorized)
}

// isCommonInternalError checks if the provided error is one of the common internal server errors defined in the commonInternalServerErrors slice.
func isCommonInternalError(e error) bool {
	for _, candidate := range commonInternalServerErrors {
		if errors.Is(e, candidate) {
			return true
		}
	}
	return false
}
