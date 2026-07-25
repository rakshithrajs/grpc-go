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

func logPrefix(fn string) string { return "[" + fn + "]: " }

const (
	fnUploadFile        = "UploadFile"
	fnDownloadFile      = "DownloadFile"
	fnListFiles         = "ListFiles"
	fnRenameFile        = "RenameFile"
	fnDeleteFile        = "DeleteFile"
	fnUserIDFromContext = "UserIDFromContext"
)

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

func NewFileHandler(fileService storage.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}
