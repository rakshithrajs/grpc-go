package handlers

import (
	filespb "cloud/gen/files/v1"
	"cloud/internal/storage"
)

type FileHandler struct {
	filespb.UnimplementedFilesServer
	storage storage.FileService
}

func NewFileHandler(fileService storage.FileService) *FileHandler {
	return &FileHandler{
		storage: fileService,
	}
}
