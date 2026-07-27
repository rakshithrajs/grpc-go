package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) DeleteFile(ctx context.Context, req *MMSpb.DeleteFileRequest) (*MMSpb.EmptyMessage, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	file, err := f.fileService.DeleteFile(ctx, req.GetFileID(), userID)
	if err != nil {
		slog.Error(logPrefix(fnDeleteFile)+"failed to delete file record", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if file == nil {
		return &MMSpb.EmptyMessage{}, nil
	}

	if err := os.Remove(file.Path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error(logPrefix(fnDeleteFile)+"failed to remove file from disk", slog.Any(config.ErrorKey, err), slog.String("path", file.Path))
		}
		if _, rbErr := f.fileService.UploadFile(ctx, file); rbErr != nil {
			slog.Error(logPrefix(fnDeleteFile)+"failed to roll back file deletion", slog.Any(config.ErrorKey, rbErr), slog.String("path", file.Path))
			return nil, status.Error(codes.Internal, ErrFailedToRollback.Error())
		}
	}

	return &MMSpb.EmptyMessage{}, nil
}
