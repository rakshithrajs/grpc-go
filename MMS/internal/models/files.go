package models

import "time"

// File represents a file owned by a user in the media storage service.
type File struct {
	ID           string    `json:"ID"`
	UserID       string    `json:"userID"`
	Name         string    `json:"name"`
	Path         string    `json:"path"`
	Size         int64     `json:"size"`
	MimeType     MimeType  `json:"mimeType"`
	CreatedAtUTC time.Time `json:"createdAtUTC"`
	UpdatedAtUTC time.Time `json:"updatedAtUTC"`
}

// UploadFileRequest contains the filename and contents for a file upload.
type UploadFileRequest struct {
	Name     string `json:"name"`
	Contents []byte `json:"contents"`
}

// RenameFileRequest contains the new filename for a file rename.
type RenameFileRequest struct {
	Name string `json:"name"`
}

// UpdateFileRequest contains optional fields for renaming or moving a file.
type UpdateFileRequest struct {
	Name string `json:"name,omitempty"`
	Path string `json:"path,omitempty"`
}

// ListFileResponse contains the metadata returned when listing a user's files.
type ListFileResponse struct {
	ID       string   `json:"ID"`
	FileName string   `json:"fileName"`
	FileSize int64    `json:"fileSize"`
	MimeType MimeType `json:"mimeType"`
}
