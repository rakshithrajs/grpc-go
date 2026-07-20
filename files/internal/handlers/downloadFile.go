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

var (
	ErrFileIDRequired       = errors.New("file ID is required")
	ErrFailedToDownloadFile = errors.New("failed to download file")
)

func (f *FileHandler) DownloadFile(ctx context.Context, req *filespb.DownloadFileRequest) (*filespb.DownloadFileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFileID() == "" {
		return nil, status.Error(codes.InvalidArgument, ErrFileIDRequired.Error())
	}

	file, err := f.fileService.GetFileByID(ctx, req.GetFileID(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return &filespb.DownloadFileResponse{}, nil
		}
		slog.Error("[DownloadFile]: failed to get file", slog.Any("error", err))
		return nil, status.Error(codes.Internal, ErrFailedToDownloadFile.Error())
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
		MimeType: *file.MimeType,
		Content:  contents,
	}, nil
}
