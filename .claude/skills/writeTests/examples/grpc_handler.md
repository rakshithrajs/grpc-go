# example only - grpc handler

```go
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
```

```go
package handlers

import (
	"context"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/mocks"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestDeleteFile(t *testing.T) {
	ctxWithUser := func() context.Context {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs(config.UserIDMetadataKey, "user-123"))
	}

	tests := []struct {
		name         string
		setupCtx     func() context.Context
		fileID       string
		mockDbErr    mocks.DbOperationError
		expectedCode codes.Code
		expectedErr  string
		expectedData *MMSpb.EmptyMessage
	}{
		{
			name: "deletion fails due to missing metadata",
			setupCtx: func() context.Context {
				return context.Background()
			},
			fileID:       "file-id-123",
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingMetadata.Error(),
		},
		{
			name: "deletion fails due to missing userID",
			setupCtx: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs())
			},
			fileID:       "file-id-123",
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingUserID.Error(),
		},
		{
			name:         "deletion succeeds with file not found",
			setupCtx:     ctxWithUser,
			fileID:       "file-id-123",
			mockDbErr:    mocks.DbOpNotFound,
			expectedErr:  config.NullString,
			expectedCode: codes.OK,
		},
		{
			name:         "deletion fails due to db internal error",
			setupCtx:     ctxWithUser,
			fileID:       "file-id-123",
			mockDbErr:    mocks.DbOpInternalError,
			expectedErr:  storage.ErrFailedToDeleteFile.Error(),
			expectedCode: codes.Internal,
		},
		{
			name:         "deletion succeeds",
			setupCtx:     ctxWithUser,
			fileID:       "file-id-123",
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
			expectedData: &MMSpb.EmptyMessage{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			mockService := &mocks.MockFileService{MockErr: tt.mockDbErr}
			handler := &FileHandler{fileService: mockService}

			req := &MMSpb.DeleteFileRequest{FileID: tt.fileID}

			// act
			resp, err := handler.DeleteFile(tt.setupCtx(), req)

			// assert
			if tt.expectedCode != status.Code(err) {
				t.Fatalf("expected code %v, got %v", tt.expectedCode, status.Code(err))
			}

			if tt.expectedErr != config.NullString && status.Convert(err).Message() != tt.expectedErr {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, status.Convert(err).Message())
			}

			if tt.expectedData != nil {
				if resp == nil {
					t.Fatalf("expected response %v, got nil", tt.expectedData)
				}
			}
		})
	}
}
```
