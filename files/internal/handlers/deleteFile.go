package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"

	filespb "github.com/rakshithrajs/cloud/services/files/gen/files/v1"
	"github.com/rakshithrajs/cloud/services/files/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) DeleteFile(ctx context.Context, req *filespb.DeleteFileRequest) (*filespb.DeleteFileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFileID() == nullString {
		return nil, status.Error(codes.InvalidArgument, ErrFileIDRequired.Error())
	}

	file, err := f.fileService.GetFileByID(ctx, req.GetFileID(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return &filespb.DeleteFileResponse{}, nil
		}
		slog.Error(logPrefix(fnDeleteFile)+"failed to get file", slog.Any("error", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := f.fileService.DeleteFile(ctx, req.GetFileID(), userID); err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return &filespb.DeleteFileResponse{}, nil
		}
		slog.Error(logPrefix(fnDeleteFile)+"failed to delete file record", slog.Any("error", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	if err := os.Remove(*file.Path); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			slog.Error(logPrefix(fnDeleteFile)+"failed to remove file from disk", slog.Any("error", err), slog.String("path", *file.Path))
		}
	}

	return &filespb.DeleteFileResponse{}, nil
}
