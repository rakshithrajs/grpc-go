package handlers

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/mocks"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestDownloadFile(t *testing.T) {
	tempDir := t.TempDir()
	config.SetConfig(&config.Config{UserStoragePath: tempDir})

	const (
		userID   = "user-123"
		fileName = "test.txt"
		fileID   = "file-id-123"
	)

	userDir := filepath.Join(tempDir, userID)
	filePath := filepath.Join(userDir, fileName)

	resetFile := func() {
		t.Helper()
		if err := os.RemoveAll(userDir); err != nil {
			t.Fatalf("failed to remove user dir: %v", err)
		}
		if err := os.MkdirAll(userDir, 0o755); err != nil {
			t.Fatalf("failed to create user dir: %v", err)
		}
		if err := os.WriteFile(filePath, []byte("test content"), 0o644); err != nil {
			t.Fatalf("failed to write test file: %v", err)
		}
	}

	ctxWithUser := func() context.Context {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs("userID", userID))
	}

	tests := []struct {
		name        string
		setupCtx    func() context.Context
		fileID      string
		mockDbErr   mocks.DbOperationError
		beforeCall  func()
		expectedCode codes.Code
		expectedErr string
		expectedData MMSpb.DownloadFileResponse
	}{
		{
			name: "download fails due to missing metadata",
			setupCtx: func() context.Context {
				return context.Background()
			},
			fileID:       fileID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingMetadata.Error(),
		},
		{
			name: "download fails due to missing userID",
			setupCtx: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs())
			},
			fileID:       fileID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingUserID.Error(),
		},
		{
			name:         "download succeeds but file not found",
			setupCtx:     ctxWithUser,
			fileID:       fileID,
			mockDbErr:    mocks.DbOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
			expectedData: MMSpb.DownloadFileResponse{},
		},
		{
			name:         "download fails due to db internal error",
			setupCtx:     ctxWithUser,
			fileID:       fileID,
			mockDbErr:    mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  ErrFailedToDownloadFile.Error(),
		},
		{
			name:        "download fails when file cannot be opened",
			setupCtx:    ctxWithUser,
			fileID:      fileID,
			beforeCall: func() { _ = os.Remove(filePath) },
			expectedCode: codes.Internal,
			expectedErr: ErrFailedToDownloadFile.Error(),
		},
		{
			name:     "download fails when file cannot be read",
			setupCtx: ctxWithUser,
			fileID:   fileID,
			beforeCall: func() {
				_ = os.Remove(filePath)
				if err := os.Mkdir(filePath, 0o755); err != nil {
					t.Fatalf("failed to create directory for read failure: %v", err)
				}
			},
			expectedCode: codes.Internal,
			expectedErr:  ErrFailedToDownloadFile.Error(),
		},
		{
			name:         "download succeeds",
			setupCtx:     ctxWithUser,
			fileID:       fileID,
			expectedCode: codes.OK,
			expectedErr:  "",
			expectedData: MMSpb.DownloadFileResponse{
				FileName: fileName,
				MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
				Content:  []byte("test content"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetFile()
			if tt.beforeCall != nil {
				tt.beforeCall()
			}

			mockService := &mocks.MockFileService{MockErr: tt.mockDbErr}
			handler := &FileHandler{fileService: mockService}

			resp, err := handler.DownloadFile(tt.setupCtx(), &MMSpb.DownloadFileRequest{FileID: tt.fileID})

			st := status.Convert(err)
			if tt.expectedCode != st.Code() {
				t.Fatalf("expected code %v, got %v", tt.expectedCode, st.Code())
			}

			if st.Message() != tt.expectedErr {
				t.Fatalf("expected error %q, got %q", tt.expectedErr, st.Message())
			}

			if tt.expectedCode != codes.OK {
				return
			}

			if resp.FileName != tt.expectedData.FileName {
				t.Errorf("expected filename %q, got %q", tt.expectedData.FileName, resp.FileName)
			}
			if resp.MimeType != tt.expectedData.MimeType {
				t.Errorf("expected mime type %v, got %v", tt.expectedData.MimeType, resp.MimeType)
			}
			if string(resp.Content) != string(tt.expectedData.Content) {
				t.Errorf("expected content %q, got %q", tt.expectedData.Content, resp.Content)
			}
		})
	}
}
