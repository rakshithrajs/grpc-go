package handlers

import (
	"context"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/mocks"
	"github.com/rakshithrajs/cloud/MMS/internal/models"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestListFiles(t *testing.T) {
	ctxWithUser := func() context.Context {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs(config.UserIDMetadataKey, "user-123"))
	}

	tests := []struct {
		name          string
		setupCtx      func() context.Context
		mockDbErr     mocks.DbOperationError
		files         []*models.ListFileResponse
		expectedCode  codes.Code
		expectedErr   string
		expectedFiles []*MMSpb.File
	}{
		{
			name: "list files fails due to missing metadata",
			setupCtx: func() context.Context {
				return context.Background()
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingMetadata.Error(),
		},
		{
			name: "list files fails due to missing userID",
			setupCtx: func() context.Context {
				return metadata.NewIncomingContext(context.Background(), metadata.Pairs())
			},
			expectedCode: codes.Unauthenticated,
			expectedErr:  ErrMissingUserID.Error(),
		},
		{
			name:         "list files fails due to db internal error",
			setupCtx:     ctxWithUser,
			mockDbErr:    mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  ErrFailedToListFiles.Error(),
		},
		{
			name:         "list files succeeds",
			setupCtx:     ctxWithUser,
			expectedCode: codes.OK,
			expectedFiles: []*MMSpb.File{
				{ID: "file-1", FileName: "file1.txt", FileSize: 100, MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN},
				{ID: "file-2", FileName: "file2.txt", FileSize: 200, MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN},
			},
		},
		{
			name:          "list files succeeds with empty list",
			setupCtx:      ctxWithUser,
			files:         []*models.ListFileResponse{},
			expectedCode:  codes.OK,
			expectedFiles: []*MMSpb.File{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			handler := NewFileHandler(&mocks.MockFileService{MockErr: tt.mockDbErr, Files: tt.files})

			// act
			resp, err := handler.ListFiles(tt.setupCtx(), &MMSpb.EmptyMessage{})

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

			gotFiles := resp.GetFile()
			if len(gotFiles) != len(tt.expectedFiles) {
				t.Errorf("expected %d files, got %d", len(tt.expectedFiles), len(gotFiles))
				return
			}

			for i, expected := range tt.expectedFiles {
				got := gotFiles[i]
				if got.GetID() != expected.GetID() {
					t.Errorf("expected file %d id %q, got %q", i, expected.GetID(), got.GetID())
				}
				if got.GetFileName() != expected.GetFileName() {
					t.Errorf("expected file %d filename %q, got %q", i, expected.GetFileName(), got.GetFileName())
				}
				if got.GetFileSize() != expected.GetFileSize() {
					t.Errorf("expected file %d size %d, got %d", i, expected.GetFileSize(), got.GetFileSize())
				}
				if got.GetMimeType() != expected.GetMimeType() {
					t.Errorf("expected file %d mime type %v, got %v", i, expected.GetMimeType(), got.GetMimeType())
				}
			}
		})
	}
}
