package handlers

import (
	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"
)

const (
	// function name for UploadFile
	fnUploadFile = "UploadFile"

	// function name for DownloadFile
	fnDownloadFile = "DownloadFile"

	// function name for RenameFile
	fnRenameFile = "RenameFile"

	// function name for DeleteFile
	fnDeleteFile = "DeleteFile"
)

// FileHandler implements the MMS gRPC file service handlers.
type FileHandler struct {
	MMSpb.UnimplementedFilesServer
	fileService storage.FileService
}

// NewFileHandler creates a new instance of FileHandler with the provided file service.
func NewFileHandler(fileService storage.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}
