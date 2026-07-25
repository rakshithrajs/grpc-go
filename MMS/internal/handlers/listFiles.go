package handlers

import (
	"context"
	"log/slog"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) ListFiles(ctx context.Context, req *MMSpb.EmptyMessage) (*MMSpb.ListFilesResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	files, err := f.fileService.GetFiles(ctx, userID)
	if err != nil {
		slog.Error(logPrefix(fnListFiles)+"failed to get files", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, ErrFailedToListFiles.Error())
	}

	respFiles := make([]*MMSpb.File, 0, len(files))
	for _, file := range files {
		respFile := &MMSpb.File{
			ID:       file.ID,
			FileName: file.FileName,
			FileSize: file.FileSize,
			MimeType: toProtoMimeType(file.MimeType),
		}
		respFiles = append(respFiles, respFile)
	}

	return &MMSpb.ListFilesResponse{
		File: respFiles,
	}, nil
}
