package handlers

import (
	filespb "cloud/gen/files/v1"
	"cloud/internal/handlers"
	"cloud/internal/storage"
	"context"
	"errors"
	"log/slog"
	"os"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) DeleteFile(ctx context.Context, req *filespb.DeleteFileRequest) (*filespb.DeleteFileResponse, error) {
	userID, err := handlers.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFileID() == "" {
		return nil, status.Error(codes.InvalidArgument, ErrFileIDRequired.Error())
	}

	file, err := f.storage.GetFileByID(ctx, req.GetFileID(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileIDDoesntExist) {
			return &filespb.DeleteFileResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := f.storage.DeleteFile(ctx, req.GetFileID(), userID); err != nil {
		if errors.Is(err, storage.ErrFileIDDoesntExist) {
			return &filespb.DeleteFileResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := os.Remove(*file.Path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error("[DeleteFile]: failed to remove file from disk", slog.Any("error", err), slog.String("path", *file.Path))
		}
	}

	return &filespb.DeleteFileResponse{}, nil
}
