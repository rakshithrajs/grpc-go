package storage

import "errors"

var (
	// ErrFailedToUploadFile is returned when a file upload fails.
	ErrFailedToUploadFile = errors.New("failed to upload file")

	// ErrFileAlreadyExists is returned when a file with the same name already exists.
	ErrFileAlreadyExists = errors.New("file already exists")

	// ErrFailedToGetFileByID is returned when fetching a file by ID fails.
	ErrFailedToGetFileByID = errors.New("failed to get file by ID")

	// ErrFailedToUpdateFile is returned when a file update fails.
	ErrFailedToUpdateFile = errors.New("failed to update file")

	// ErrFailedToDeleteFile is returned when a file deletion fails.
	ErrFailedToDeleteFile = errors.New("failed to delete file")

	// ErrFailedToRollback is returned when a rollback operation fails.
	ErrFailedToRollback = errors.New("failed to roll back")
)
