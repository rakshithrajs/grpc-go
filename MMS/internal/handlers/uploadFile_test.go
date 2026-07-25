package handlers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/MMS/internal/mocks/utils"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUploadFile(t *testing.T) {
	tempDir := t.TempDir()
	config.SetConfig(&config.Config{UserStoragePath: tempDir})

	const userID = "user-123"

	ctxWithUser := func() context.Context {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs(config.UserIDMetadataKey, userID))
	}

	userDir := filepath.Join(tempDir, userID)

	resetDisk := func() {
		t.Helper()
		if err := os.RemoveAll(userDir); err != nil {
			t.Fatalf("failed to remove user dir: %v", err)
		}
	}

	tests := []struct {
		name         string
		setupCtx     func() context.Context
		fileName     string
		content      []byte
		mockDbErr    mocks.DbOperationError
		preCreate    bool
		expectedCode codes.Code
		expectedErr  string
		expectedFile *MMSpb.File
	}{
		{
			name: "upload fails due to missing metadata",
			setupCtx: func() context.Context {
				return context.Background()
			},
			fileName:     "missing-metadata.txt",
			content:      []byte("content"),
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingMetadata.Error(),
		},
		{
			name: "upload fails due to missing userID",
			setupCtx: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs())
			},
			fileName:     "missing-userid.txt",
			content:      []byte("content"),
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingUserID.Error(),
		},
		{
			name:         "upload fails as file path already exists",
			setupCtx:     ctxWithUser,
			fileName:     "already-exists.txt",
			content:      []byte("content"),
			preCreate:    true,
			expectedCode: codes.AlreadyExists,
			expectedErr:  storage.ErrFilePathAlreadyExists.Error(),
		},
		{
			name:         "upload fails as file name already exists in db",
			setupCtx:     ctxWithUser,
			fileName:     "db-duplicate.txt",
			content:      []byte("content"),
			mockDbErr:    mocks.DbOpDuplicateName,
			expectedCode: codes.AlreadyExists,
			expectedErr:  storage.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "upload fails due to db internal error",
			setupCtx:     ctxWithUser,
			fileName:     "db-internal.txt",
			content:      []byte("content"),
			mockDbErr:    mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload succeeds",
			setupCtx:     ctxWithUser,
			fileName:     "success.txt",
			content:      []byte("content"),
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
			expectedFile: &MMSpb.File{
				ID:       "file-id-123",
				FileName: "success.txt",
				FileSize: 7,
				MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			resetDisk()
			if tt.preCreate {
				if err := os.MkdirAll(userDir, 0o755); err != nil {
					t.Fatalf("failed to create user dir: %v", err)
				}
				if err := os.WriteFile(filepath.Join(userDir, tt.fileName), []byte("existing"), 0o644); err != nil {
					t.Fatalf("failed to create existing file: %v", err)
				}
			}

			handler := NewFileHandler(&mocks.MockFileService{MockErr: tt.mockDbErr})

			// act
			resp, err := handler.UploadFile(tt.setupCtx(), &MMSpb.UploadFileRequest{
				FileName: tt.fileName,
				Content:  tt.content,
			})

			// assert
			st := status.Convert(err)
			if tt.expectedCode != st.Code() {
				t.Errorf("expected code %v, got %v", tt.expectedCode, st.Code())
			}

			if st.Message() != tt.expectedErr {
				t.Errorf("expected error %q, got %q", tt.expectedErr, st.Message())
			}

			if tt.expectedCode == codes.OK {
				mockUtils.CheckData(t, resp, &MMSpb.UploadFileResponse{File: tt.expectedFile})
			}
		})
	}
}
