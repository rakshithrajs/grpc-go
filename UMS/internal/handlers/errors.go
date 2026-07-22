package handlers

import "errors"

var (
	ErrInvalidJSON          = errors.New("invalid JSON payload")
	ErrFailedToRegisterUser = errors.New("failed to register user")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrFailedToLoginUser    = errors.New("failed to login user")
	ErrNoFieldsToUpdate     = errors.New("no fields to update")
	ErrInvalidID            = errors.New("invalid ID")
	ErrFileIDRequired       = errors.New("file ID is required")
	ErrFileIsRequired       = errors.New("file is required")
	ErrFailedToUploadFile   = errors.New("failed to upload file")
	ErrFailedToRenameFile   = errors.New("failed to rename file")
	ErrFailedToListFiles    = errors.New("failed to list files")
	ErrFailedToDownloadFile = errors.New("failed to download file")
	ErrFailedToDeleteFile   = errors.New("failed to delete file")
	ErrFileRenamed          = errors.New("file renamed successfully")
	ErrFileDeleted          = errors.New("file deleted successfully")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrFileNotFound         = errors.New("file not found")
	ErrSomethingWentWrong   = errors.New("something went wrong")
)
