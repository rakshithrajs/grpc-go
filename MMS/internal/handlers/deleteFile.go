package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedToDeleteFile = errors.New("failed to delete file")
)

func (f *FileHandler) DeleteFile(ctx context.Context, req *MMSpb.DeleteFileRequest) (*MMSpb.EmptyMessage, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	file, err := f.fileService.DeleteFile(ctx, req.GetFileID(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return &MMSpb.EmptyMessage{}, nil
		}
		slog.Error(logPrefix(fnDeleteFile)+"failed to delete file record", slog.Any("error", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := os.Remove(file.Path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error(logPrefix(fnDeleteFile)+"failed to remove file from disk", slog.Any("error", err), slog.String("path", file.Path))
		}
	}

	return &MMSpb.EmptyMessage{}, nil
}
