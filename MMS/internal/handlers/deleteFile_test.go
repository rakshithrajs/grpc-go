package handlers

import (
	"context"
	"os"
	"path/filepath"
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
	tempDir := t.TempDir()
	config.SetConfig(&config.Config{UserStoragePath: tempDir})

	const (
		userID   = "user-123"
		fileID   = "file-id-123"
		fileName = "test.txt"
	)

	ctxWithUser := func() context.Context {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs(config.UserIDMetadataKey, userID))
	}

	filePath := filepath.Join(tempDir, userID, fileName)

	reset := func() {
		t.Helper()
		_ = os.RemoveAll(filepath.Join(tempDir, userID))
	}

	tests := []struct {
		name            string
		setupCtx        func() context.Context
		mockDeleteErr   mocks.DbOperationError
		mockRollbackErr mocks.DbOperationError
		beforeCall      func()
		expectedCode    codes.Code
		expectedErr     string
	}{
		{
			name: "deletion fails due to missing metadata",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingMetadata.Error(),
		},
		{
			name: "deletion fails due to missing userID",
			setupCtx: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs())
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingUserID.Error(),
		},
		{
			name:          "deletion succeeds with file not found",
			setupCtx:      ctxWithUser,
			mockDeleteErr: mocks.DbOpNotFound,
			expectedCode:  codes.OK,
			expectedErr:   config.NullString,
		},
		{
			name:          "deletion fails due to db internal error",
			setupCtx:      ctxWithUser,
			mockDeleteErr: mocks.DbOpInternalError,
			expectedCode:  codes.Internal,
			expectedErr:   storage.ErrFailedToDeleteFile.Error(),
		},
		{
			name:            "deletion fails due to rollback failure",
			setupCtx:        ctxWithUser,
			mockRollbackErr: mocks.DbOpInternalError,
			beforeCall: func() {
				if err := os.MkdirAll(filePath, 0o755); err != nil {
					t.Fatalf("failed to create dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(filePath, "blocker"), []byte("x"), 0o644); err != nil {
					t.Fatalf("failed to write blocker: %v", err)
				}
			},
			expectedCode: codes.Internal,
			expectedErr:  ErrFailedToRollback.Error(),
		},
		{
			name: "deletion succeeds when disk file is already gone",
			setupCtx: ctxWithUser,
			beforeCall: func() {
				if err := os.RemoveAll(filePath); err != nil {
					t.Fatalf("failed to remove file before delete call: %v", err)
				}
			},
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "deletion succeeds",
			setupCtx:     ctxWithUser,
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reset()
			if tt.beforeCall != nil {
				tt.beforeCall()
			}

			mockService := &mocks.MockFileService{
				DeleteFileErr: tt.mockDeleteErr,
				UploadFileErr: tt.mockRollbackErr,
			}
			handler := &FileHandler{fileService: mockService}

			resp, err := handler.DeleteFile(tt.setupCtx(), &MMSpb.DeleteFileRequest{FileID: fileID})

			if tt.expectedCode != status.Code(err) {
				t.Fatalf("expected code %v, got %v", tt.expectedCode, status.Code(err))
			}

			if tt.expectedErr != config.NullString && status.Convert(err).Message() != tt.expectedErr {
				t.Fatalf("expected error %v, got %v", tt.expectedErr, status.Convert(err).Message())
			}

			if tt.expectedCode == codes.OK && resp == nil {
				t.Fatalf("expected response, got nil")
			}
		})
	}
}
