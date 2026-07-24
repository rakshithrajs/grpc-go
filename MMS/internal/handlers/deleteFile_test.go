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
	}{
		{
			name: "missing metadata",
			setupCtx: func() context.Context {
				return context.Background()
			},
			fileID:       "file-id-123",
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingMetadata.Error(),
		},
		{
			name: "missing userID",
			setupCtx: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs())
			},
			fileID:       "file-id-123",
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingUserID.Error(),
		},
		{
			name:         "file not found",
			setupCtx:     ctxWithUser,
			fileID:       "file-id-123",
			mockDbErr:    mocks.DbOpNotFound,
			expectedErr:  "",
			expectedCode: codes.OK,
		},
		{
			name:         "db internal error",
			setupCtx:     ctxWithUser,
			fileID:       "file-id-123",
			mockDbErr:    mocks.DbOpInternalError,
			expectedErr:  storage.ErrFailedToDeleteFile.Error(),
			expectedCode: codes.Internal,
		},
		{
			name:         "delete succeeds",
			setupCtx:     ctxWithUser,
			fileID:       "file-id-123",
			expectedCode: codes.OK,
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			mockService := &mocks.MockFileService{MockErr: tt.mockDbErr}
			handler := &FileHandler{fileService: mockService}

			req := &MMSpb.DeleteFileRequest{FileID: tt.fileID}

			// act
			_, err := handler.DeleteFile(tt.setupCtx(), req)

			// assert
			if tt.expectedCode != status.Code(err) {
				t.Fatalf("expected code %v, got %v", tt.expectedCode, status.Code(err))
			}

			if tt.expectedErr != "" && status.Convert(err).Message() != tt.expectedErr {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, status.Convert(err).Message())
			}
		})
	}
}
