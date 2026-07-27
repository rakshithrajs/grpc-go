package handlers

import (
	"errors"
)

var (
	// message: failed to download file
	ErrFailedToDownloadFile = errors.New("failed to download file")

	// message: missing metadata in context
	ErrMissingMetadata = errors.New("missing metadata")

	// message: missing user ID in metadata
	ErrMissingUserID = errors.New("missing user ID in metadata")

	// message: failed to upload file
	ErrFailedToListFiles = errors.New("failed to list files")

	// message: failed to rename file
	ErrFailedToRenameFile = errors.New("failed to rename file")

	// message: failed to roll back
	ErrFailedToRollback = errors.New("failed to roll back")
)
