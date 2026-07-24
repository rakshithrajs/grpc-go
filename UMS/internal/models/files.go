package models

type UserFiles struct {
	FileID   string `json:"fileId"`
	UserID   string `json:"userId"`
	FileName string `json:"fileName"`
}

type File struct {
	ID       string `json:"id"`
	FileName string `json:"fileName"`
	FileSize int64  `json:"fileSize"`
	MimeType string `json:"mimeType"`
}

type RenameFileRequest struct {
	NewName string `json:"newName" validate:"required,isValueEmpty,max=255"`
}

type UploadRequest struct {
	UserID   string `validate:"required,isValueEmpty"`
	FileName string `validate:"required,isValueEmpty"`
}
