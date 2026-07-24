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

func TestRenameFile(t *testing.T) {
	tempDir := t.TempDir()
	config.SetConfig(&config.Config{UserStoragePath: tempDir})

	const (
		userID  = "user-123"
		fileID  = "file-id-123"
		newName = "renamed.txt"
		oldName = "old.txt"
	)

	ctxWithUser := func() context.Context {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs("userID", userID))
	}

	userDir := filepath.Join(tempDir, userID)
	oldPath := filepath.Join(userDir, oldName)
	newPath := filepath.Join(userDir, newName)

	resetDisk := func() {
		t.Helper()
		if err := os.RemoveAll(userDir); err != nil {
			t.Fatalf("failed to remove user dir: %v", err)
		}
		if err := os.MkdirAll(userDir, 0o755); err != nil {
			t.Fatalf("failed to create user dir: %v", err)
		}
	}

	tests := []struct {
		name          string
		setupCtx      func() context.Context
		mockDbErr     mocks.DbOperationError
		returnOldName bool
		preCreate     bool
		expectedCode  codes.Code
		expectedErr   string
	}{
		{
			name: "rename fails due to missing metadata",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingMetadata.Error(),
		},
		{
			name: "rename fails due to missing userID",
			setupCtx: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs())
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingUserID.Error(),
		},
		{
			name:         "rename fails due to db internal error",
			setupCtx:     ctxWithUser,
			mockDbErr:    mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  ErrFailedToRenameFile.Error(),
		},
		{
			name:         "rename succeeds but file not found",
			setupCtx:     ctxWithUser,
			mockDbErr:    mocks.DbOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:         "rename fails as file name already exists",
			setupCtx:     ctxWithUser,
			mockDbErr:    mocks.DbOpDuplicateName,
			expectedCode: codes.AlreadyExists,
			expectedErr:  storage.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "rename succeeds early when name is unchanged",
			setupCtx:     ctxWithUser,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:          "rename succeeds with disk rename",
			setupCtx:      ctxWithUser,
			returnOldName: true,
			preCreate:     true,
			expectedCode:  codes.OK,
			expectedErr:   "",
		},
		{
			name:          "rename fails and rolls back disk rename",
			setupCtx:      ctxWithUser,
			returnOldName: true,
			expectedCode:  codes.Internal,
			expectedErr:   ErrFailedToRenameFile.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			resetDisk()
			if tt.preCreate {
				if err := os.WriteFile(oldPath, []byte("content"), 0o644); err != nil {
					t.Fatalf("failed to create old file: %v", err)
				}
			}

			mock := &mocks.MockFileService{
				MockErr:       tt.mockDbErr,
				ReturnOldName: tt.returnOldName,
			}
			handler := NewFileHandler(mock)

			// act
			_, err := handler.RenameFile(tt.setupCtx(), &MMSpb.RenameFileRequest{
				FileID:  fileID,
				NewName: newName,
			})

			// assert
			st := status.Convert(err)
			if tt.expectedCode != st.Code() {
				t.Errorf("expected code %v, got %v", tt.expectedCode, st.Code())
			}

			if st.Message() != tt.expectedErr {
				t.Errorf("expected error %q, got %q", tt.expectedErr, st.Message())
			}

			if tt.expectedCode != codes.OK {
				return
			}

			if tt.name == "rename succeeds with disk rename" {
				if _, statErr := os.Stat(oldPath); !os.IsNotExist(statErr) {
					t.Errorf("expected old file to be removed, still exists at %s", oldPath)
				}
				if _, statErr := os.Stat(newPath); os.IsNotExist(statErr) {
					t.Errorf("expected new file to exist, missing at %s", newPath)
				}
			}
		})
	}
}
