package handlers

import (
	"errors"
)

var (
	ErrFailedToDownloadFile = errors.New("failed to download file")
	ErrMissingMetadata      = errors.New("missing metadata")
	ErrMissingUserID        = errors.New("missing user id in metadata")
	ErrFailedToListFiles    = errors.New("failed to list files")
	ErrFailedToRenameFile   = errors.New("failed to rename file")
)
