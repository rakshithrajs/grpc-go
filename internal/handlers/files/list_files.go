package handlers

import (
	filespb "cloud/gen/files/v1"
	"cloud/internal/handlers"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) ListFiles(ctx context.Context, req *filespb.ListFilesRequest) (*filespb.ListFilesResponse, error) {
	userID, err := handlers.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	files, err := f.storage.GetFiles(ctx, userID)
	if err != nil {
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
