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

var (
	ErrFileIDRequired       = errors.New("file ID is required")
	ErrFailedToDownloadFile = errors.New("failed to download file")
)

func (f *FileHandler) DownloadFile(ctx context.Context, req *filespb.DownloadFileRequest) (*filespb.DownloadFileResponse, error) {
	userID, err := handlers.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFileID() == "" {
		return nil, status.Error(codes.InvalidArgument, ErrFileIDRequired.Error())
	}

	file, err := f.storage.GetFileByID(ctx, req.FileID, userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileIDDoesntExist) {
			return &filespb.DownloadFileResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	fi, err := os.Open(*file.Path)
	if err != nil {
		slog.Error("[DownloadFile]: failed to open file", slog.Any("error", err), slog.String("path", *file.Path))
		return nil, status.Error(codes.Internal, ErrFailedToDownloadFile.Error())
	}
	defer fi.Close()

	contents := make([]byte, *file.Size)
	if _, err := fi.Read(contents); err != nil {
		slog.Error("[DownloadFile]: failed to read file", slog.Any("error", err), slog.String("path", *file.Path))
		return nil, status.Error(codes.Internal, ErrFailedToDownloadFile.Error())
	}

	return &filespb.DownloadFileResponse{
		FileName: *file.Name,
		Content:  contents,
		MimeType: *file.MimeType,
	}, nil
}
