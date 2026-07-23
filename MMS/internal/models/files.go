package models

import "time"

type File struct {
	ID           *string    `json:"id"`
	UserID       *string    `json:"userID"`
	Name         *string    `json:"name"`
	Path         *string    `json:"path"`
	Size         *int64     `json:"size"`
	MimeType     *MimeType  `json:"mimeType"`
	CreatedAtUTC *time.Time `json:"createdAtUTC"`
	UpdatedAtUTC *time.Time `json:"updatedAtUTC"`
}

type UploadFileRequest struct {
	Name     *string `json:"name" validate:"required,isValueEmpty,isValidFileName,max=150"`
	Contents []byte  `json:"contents" validate:"required"`
}

type RenameFileRequest struct {
	Name *string `json:"name" validate:"required,isValueEmpty,isValidFileName,max=150"`
}

type UpdateFileRequest struct {
	Name *string `json:"name,omitempty"`
	Path *string `json:"path,omitempty"`
}

type ListFileResponse struct {
	ID       *string   `json:"id"`
	FileName *string   `json:"fileName"`
	FileSize *int64    `json:"fileSize"`
	MimeType *MimeType `json:"mimeType"`
}
