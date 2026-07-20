package models

import "time"

type File struct {
	ID           *string    `json:"id"`
	UserID       *string    `json:"userID"`
	Name         *string    `json:"name"`
	Path         *string    `json:"path"`
	Size         *int64     `json:"size"`
	MimeType     *string    `json:"mimeType"`
	CreatedAtUTC *time.Time `json:"createdAtUTC"`
	UpdatedAtUTC *time.Time `json:"updatedAtUTC"`
}

type UploadFileRequest struct {
	Name     *string `validate:"required,isValueEmpty,isValidFileName,max=150"`
	Contents []byte  `validate:"required"`
}

type RenameFileRequest struct {
	Name *string `validate:"required,isValueEmpty,isValidFileName,max=150"`
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
