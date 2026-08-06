package errors

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

// Source identifiers used when mapping errors to HTTP responses.
const (
	// AuthMiddlewareSource identifies errors originating from the auth middleware.
	AuthMiddlewareSource = "AuthMiddleware"

	// LoginUserSource identifies errors originating from the login handler.
	LoginUserSource = "LoginUserHandler"

	// RegisterUserSource identifies errors originating from the register handler.
	RegisterUserSource = "RegisterUserHandler"

	// UpdateUserSource identifies errors originating from the update user handler.
	UpdateUserSource = "UpdateUserProfileHandler"

	// GetUserProfileSource identifies errors originating from the get user profile handler.
	GetUserProfileSource = "GetUserProfileHandler"

	// UploadFileSource identifies errors originating from the upload file handler.
	UploadFileSource = "UploadFile"

	// DownloadFileSource identifies errors originating from the download file handler.
	DownloadFileSource = "DownloadFile"

	// ListFilesSource identifies errors originating from the list files handler.
	ListFilesSource = "ListFiles"

	// RenameFileSource identifies errors originating from the rename file handler.
	RenameFileSource = "RenameFile"

	// DeleteFileSource identifies errors originating from the delete file handler.
	DeleteFileSource = "DeleteFile"
)

var (
	// ErrSomethingWentWrong is returned for unexpected internal failures.
	ErrSomethingWentWrong = errors.New("something went wrong")

	// ErrMissingAuthHeader is returned when the Authorization header is missing.
	ErrMissingAuthHeader = errors.New("missing Authorization header")

	// ErrUnauthorized is returned when the request is not authenticated.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrInvalidJSON is returned when the request body is not valid JSON.
	ErrInvalidJSON = errors.New("invalid JSON")

	// ErrInvalidURI is returned when URI path parameters are invalid.
	ErrInvalidURI = errors.New("invalid URI")

	// ErrFailedToRegisterUser is returned when user registration fails internally.
	ErrFailedToRegisterUser = errors.New("failed to register user")

	// ErrInvalidCredentials is returned when login credentials are invalid.
	ErrInvalidCredentials = errors.New("invalid credentials")

	// ErrFailedToLoginUser is returned when user login fails internally.
	ErrFailedToLoginUser = errors.New("failed to login user")

	// ErrNoFieldsToUpdate is returned when an update request has no fields to change.
	ErrNoFieldsToUpdate = errors.New("no fields to update")

	// ErrPasswordSameAsOldPassword is returned when the new password matches the old password.
	ErrPasswordSameAsOldPassword = errors.New("new password is same as old password")

	// ErrUserEmailAlreadyExists is returned when a user with the email already exists.
	ErrUserEmailAlreadyExists = errors.New("user email already exists")

	// ErrPhoneNumberAlreadyExists is returned when a user with the phone number already exists.
	ErrPhoneNumberAlreadyExists = errors.New("phone number already exists")

	// ErrFileIsRequired is returned when a file is not provided.
	ErrFileIsRequired = errors.New("file is required")

	// ErrFileNameRequired is returned when a filename is not provided.
	ErrFileNameRequired = errors.New("file name is required")

	// ErrEmptyFileContent is returned when the provided file is empty.
	ErrEmptyFileContent = errors.New("file content is empty")

	// ErrFailedToUploadFile is returned when a file upload fails internally.
	ErrFailedToUploadFile = errors.New("failed to upload file")

	// ErrFailedToRenameFile is returned when a file rename fails internally.
	ErrFailedToRenameFile = errors.New("failed to rename file")

	// ErrFailedToListFiles is returned when listing files fails internally.
	ErrFailedToListFiles = errors.New("failed to list files")

	// ErrFailedToDownloadFile is returned when a file download fails internally.
	ErrFailedToDownloadFile = errors.New("failed to download file")

	// ErrFailedToDeleteFile is returned when a file deletion fails internally.
	ErrFailedToDeleteFile = errors.New("failed to delete file")

	// ErrUserFileAlreadyExists is returned when a user-file mapping already exists.
	ErrUserFileAlreadyExists = errors.New("user file mapping already exists")

	// ErrFailedToRollback is returned when a rollback operation fails.
	ErrFailedToRollback = errors.New("failed to rollback changes")

	// ErrFailedToCreateUser is returned when creating a user fails internally.
	ErrFailedToCreateUser = errors.New("failed to create user")

	// ErrFailedToGetUserByID is returned when fetching a user by ID fails internally.
	ErrFailedToGetUserByID = errors.New("failed to get user by ID")

	// ErrFailedToGetUserByEmail is returned when fetching a user by email fails internally.
	ErrFailedToGetUserByEmail = errors.New("failed to get user by email")

	// ErrFailedToUpdateUser is returned when updating a user fails internally.
	ErrFailedToUpdateUser = errors.New("failed to update user")

	// ErrFailedToCreateUserFile is returned when creating a user-file mapping fails internally.
	ErrFailedToCreateUserFile = errors.New("failed to create user file mapping")

	// ErrFailedToDeleteUserFile is returned when deleting a user-file mapping fails internally.
	ErrFailedToDeleteUserFile = errors.New("failed to delete user file mapping")

	// ErrFailedToUpdateUserFile is returned when updating a user-file mapping fails internally.
	ErrFailedToUpdateUserFile = errors.New("failed to update user file mapping")

	// ErrFailedToListUserFiles is returned when listing user files fails internally.
	ErrFailedToListUserFiles = errors.New("failed to list user files")
)

var errorResponseSpec = map[string]map[error]int{
	AuthMiddlewareSource: {
		ErrMissingAuthHeader: http.StatusUnauthorized,
		ErrUnauthorized:      http.StatusUnauthorized,
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
		ErrUserEmailAlreadyExists:    http.StatusConflict,
		ErrPhoneNumberAlreadyExists:  http.StatusConflict,
		ErrFailedToUpdateUser:        http.StatusInternalServerError,
	},
	UploadFileSource: {
		ErrFileIsRequired:   http.StatusBadRequest,
		ErrFileNameRequired: http.StatusBadRequest,
		ErrEmptyFileContent: http.StatusBadRequest,
	},
	DownloadFileSource: {
		ErrInvalidURI:                http.StatusBadRequest,
		modelUtils.ErrFileIDRequired: http.StatusBadRequest,
	},
	ListFilesSource: {
		ErrFailedToListFiles: http.StatusInternalServerError,
	},
	RenameFileSource: {
		ErrInvalidJSON:                  http.StatusBadRequest,
		ErrInvalidURI:                   http.StatusBadRequest,
		modelUtils.ErrFileIDRequired:    http.StatusBadRequest,
		modelUtils.ErrFileIDInvalidUUID: http.StatusBadRequest,
	},
	DeleteFileSource: {
		ErrInvalidURI:                http.StatusBadRequest,
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

// ReturnErrorResponse sends an HTTP error response based on the provided error and source.
func ReturnErrorResponse(c *gin.Context, err any, source string, defaultError error) {
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

// isAuthError reports whether the error is an authentication failure.
func isAuthError(e error) bool {
	return errors.Is(e, ErrMissingAuthHeader) ||
		errors.Is(e, ErrUnauthorized)
}

// isCommonInternalError reports whether the error is a common internal server error.
func isCommonInternalError(e error) bool {
	for _, candidate := range commonInternalServerErrors {
		if errors.Is(e, candidate) {
			return true
		}
	}
	return false
}
