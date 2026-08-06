package handlers

import (
	"errors"
)

var (
	// ErrFailedToDownloadFile is returned when a file download fails.
	ErrFailedToDownloadFile = errors.New("failed to download file")

	// ErrFailedToDeleteFile is returned when a file deletion fails.
	ErrFailedToDeleteFile = errors.New("failed to delete file")

	// ErrMissingMetadata is returned when gRPC metadata is missing from the request context.
	ErrMissingMetadata = errors.New("missing metadata")

	// ErrMissingUserID is returned when the user ID is missing from gRPC metadata.
	ErrMissingUserID = errors.New("missing user ID in metadata")

	// ErrFailedToListFiles is returned when listing files fails.
	ErrFailedToListFiles = errors.New("failed to list files")

	// ErrFailedToRenameFile is returned when a file rename fails.
	ErrFailedToRenameFile = errors.New("failed to rename file")

	// ErrFailedToRollback is returned when a rollback operation fails.
	ErrFailedToRollback = errors.New("failed to roll back")
)
