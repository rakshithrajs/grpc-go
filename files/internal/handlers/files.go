package handlers

import (
	"context"
	"errors"

	filespb "github.com/rakshithrajs/cloud/services/files/gen/files/v1"
	"github.com/rakshithrajs/cloud/services/files/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingMetadata = errors.New("missing metadata")
	ErrMissingUserID   = errors.New("missing user_id in metadata")
)

func UserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 || userIDs[0] == "" {
		return "", status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	}

	return userIDs[0], nil
}

type FileHandler struct {
	filespb.UnimplementedFilesServer
	fileService storage.FileService
}

func NewFileHandler(fileService storage.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}
