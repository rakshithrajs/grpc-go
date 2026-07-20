package handlers

import (
	"context"
	"log/slog"

	filespb "github.com/rakshithrajs/cloud/services/files/gen/files/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) ListFiles(ctx context.Context, req *filespb.ListFilesRequest) (*filespb.ListFilesResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	files, err := f.fileService.GetFiles(ctx, userID)
	if err != nil {
		slog.Error(logPrefix(fnListFiles)+"failed to get files", slog.Any("error", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	respFiles := make([]*filespb.File, 0, len(files))
	for _, file := range files {
		respFile := &filespb.File{
			ID:       *file.ID,
			FileName: *file.FileName,
			FileSize: *file.FileSize,
			MimeType: *file.MimeType,
		}
		respFiles = append(respFiles, respFile)
	}

	return &filespb.ListFilesResponse{
		File: respFiles,
	}, nil
}
