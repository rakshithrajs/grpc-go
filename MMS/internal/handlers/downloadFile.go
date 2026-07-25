package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) DownloadFile(ctx context.Context, req *MMSpb.DownloadFileRequest) (*MMSpb.DownloadFileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	file, err := f.fileService.GetFileByID(ctx, req.GetFileID(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return &MMSpb.DownloadFileResponse{}, nil
		}
		slog.Error(logPrefix(fnDownloadFile)+"failed to get file", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, ErrFailedToDownloadFile.Error())
	}

	fi, err := os.Open(file.Path)
	if err != nil {
		slog.Error(logPrefix(fnDownloadFile)+"failed to open file", slog.Any(config.ErrorKey, err), slog.String("path", file.Path))
		return nil, status.Error(codes.Internal, ErrFailedToDownloadFile.Error())
	}
	defer fi.Close()

	contents := make([]byte, file.Size)
	if _, err := fi.Read(contents); err != nil {
		slog.Error(logPrefix(fnDownloadFile)+"failed to read file", slog.Any(config.ErrorKey, err), slog.String("path", file.Path))
		return nil, status.Error(codes.Internal, ErrFailedToDownloadFile.Error())
	}

	return &MMSpb.DownloadFileResponse{
		FileName: file.Name,
		MimeType: toProtoMimeType(file.MimeType),
		Content:  contents,
	}, nil
}
