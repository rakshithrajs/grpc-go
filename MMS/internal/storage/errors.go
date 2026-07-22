package storage

import "errors"

var (
	ErrFailedToUploadFile    = errors.New("failed to upload file")
	ErrFileNameAlreadyExists = errors.New("file name already exists")
	ErrFilePathAlreadyExists = errors.New("file path already exists")
	ErrFailedToGetFiles      = errors.New("failed to get files")
	ErrFailedToGetFileByID   = errors.New("failed to get file by ID")
	ErrFileNotFound          = errors.New("file not found")
	ErrFailedToUpdateFile    = errors.New("failed to update file")
	ErrFailedToDeleteFile    = errors.New("failed to delete file")
)
