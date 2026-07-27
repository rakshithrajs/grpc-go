package storage

import "errors"

var (
	// message: failed to upload file
	ErrFailedToUploadFile = errors.New("failed to upload file")

	// message: file already exists
	ErrFileAlreadyExists = errors.New("file already exists")

	// message: file not found
	ErrFileNotFound = errors.New("file not found")

	// message: failed to get file by ID
	ErrFailedToGetFileByID = errors.New("failed to get file by ID")

	// message: failed to update file
	ErrFailedToUpdateFile = errors.New("failed to update file")

	// message: failed to delete file
	ErrFailedToDeleteFile = errors.New("failed to delete file")

	// message: failed to roll back
	ErrFailedToRollback = errors.New("failed to roll back")
)
