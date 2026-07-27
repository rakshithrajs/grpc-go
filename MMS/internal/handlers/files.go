package handlers

import (
	"context"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// logPrefix returns a formatted string for logging purposes, including the function name.
func logPrefix(fn string) string { return "[" + fn + "]: " }

const (
	// function name for UploadFile
	fnUploadFile = "UploadFile"

	// function name for DownloadFile
	fnDownloadFile = "DownloadFile"

	// function name for ListFiles
	fnListFiles = "ListFiles"

	// function name for RenameFile
	fnRenameFile = "RenameFile"

	// function name for DeleteFile
	fnDeleteFile = "DeleteFile"

	// function name for UserIDFromContext
	fnUserIDFromContext = "UserIDFromContext"
)

// UserIDFromContext extracts the user ID from the gRPC context metadata.
func UserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return config.NullString, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	}

	userIDs := md.Get(config.UserIDMetadataKey)
	if len(userIDs) == 0 || userIDs[0] == config.NullString {
		return config.NullString, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	}

	return userIDs[0], nil
}

type FileHandler struct {
	MMSpb.UnimplementedFilesServer
	fileService storage.FileService
}

// NewFileHandler creates a new instance of FileHandler with the provided file service.
func NewFileHandler(fileService storage.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}
