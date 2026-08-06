package models

// UserFiles links a file ID and name to a user.
type UserFiles struct {
	FileID   string `json:"fileID"`
	UserID   string `json:"userID"`
	FileName string `json:"fileName"`
}

// File represents a file returned by the MMS gRPC service.
type File struct {
	ID       string `json:"ID"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	MimeType string `json:"mimeType"`
}

// RenameFileRequest contains the new filename for a file rename.
type RenameFileRequest struct {
	NewName string `json:"newName" validate:"required,isValueEmpty,max=255"`
}

// FileIDURI binds and validates a fileID URI path parameter.
type FileIDURI struct {
	FileID string `uri:"fileid" json:"fileID" validate:"required,isValueEmpty,uuid"`
}
