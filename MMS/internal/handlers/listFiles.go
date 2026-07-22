package handlers

import (
	"context"
	"log/slog"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) ListMMS(ctx context.Context, req *MMSpb.ListFilesRequest) (*MMSpb.ListFilesResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	MMS, err := f.MMService.GetMMS(ctx, userID)
	if err != nil {
		slog.Error(logPrefix(fnListMMS)+"failed to get MMS", slog.Any("error", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	respMMS := make([]*MMSpb.File, 0, len(MMS))
	for _, file := range MMS {
		respFile := &MMSpb.File{
			ID:       *file.ID,
			FileName: *file.FileName,
			FileSize: *file.FileSize,
			MimeType: *file.MimeType,
		}
		respMMS = append(respMMS, respFile)
	}

	return &MMSpb.ListFilesResponse{
		File: respMMS,
	}, nil
}
