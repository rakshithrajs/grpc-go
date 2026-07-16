package models

import "time"

type File struct {
	ID           *string
	UserID       *string
	Name         *string
	Path         *string
	Size         *int64
	MimeType     *string
	CreatedAtUTC *time.Time
	UpdatedAtUTC *time.Time
}

type UploadFileRequest struct {
	Name     *string `validate:"required,isValidFileName,max=150"`
	Contents []byte  `validate:"required"`
}

type RenameFileRequest struct {
	Name *string `validate:"required,isValidFileName,max=150"`
}

type UpdateFileRequest struct {
	Name *string
	Path *string
}

type ListFileResponse struct {
	ID       *string
	FileName *string
	FileSize *int64
	MimeType *string
}
