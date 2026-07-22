package handlers

import (
	"context"
	"errors"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const nullString = ""

func logPrefix(fn string) string { return "[" + fn + "]: " }

const (
	fnUploadFile        = "UploadFile"
	fnDownloadFile      = "DownloadFile"
	fnListFiles         = "ListFiles"
	fnRenameFile        = "RenameFile"
	fnDeleteFile        = "DeleteFile"
	fnUserIDFromContext = "UserIDFromContext"
)

var (
	ErrMissingMetadata = errors.New("missing metadata")
	ErrMissingUserID   = errors.New("missing user id in metadata")
)

func UserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nullString, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	}

	userIDs := md.Get("x-user-id")
	if len(userIDs) == 0 || userIDs[0] == nullString {
		return nullString, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
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
