package handlers

import (
	"context"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/mocks"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestListFiles(t *testing.T) {
	ctxWithUser := func() context.Context {
		return metadata.NewIncomingContext(context.Background(), metadata.Pairs("userID", "user-123"))
	}

	tests := []struct {
		name          string
		setupCtx      func() context.Context
		mockDbErr     mocks.DbOperationError
		expectedCode  codes.Code
		expectedErr   string
		expectedFiles int
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
			expectedErr:  storage.ErrFailedToGetFiles.Error(),
		},
		{
			name:          "list files succeeds",
			setupCtx:      ctxWithUser,
			expectedCode:  codes.OK,
			expectedFiles: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			handler := NewFileHandler(&mocks.MockFileService{MockErr: tt.mockDbErr})

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

			if len(resp.GetFile()) != tt.expectedFiles {
				t.Errorf("expected %d files, got %d", tt.expectedFiles, len(resp.GetFile()))
			}
		})
	}
}
