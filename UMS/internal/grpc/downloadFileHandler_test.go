package grpc

import (
	"context"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDownloadFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name         string
		mockGrpcErr  mocks.GrpcOperationError
		expectedCode codes.Code
		expectedErr  string
		expectedData MMSpb.DownloadFileResponse
	}{
		{
			name:         "file download fails as missing metadata",
			mockGrpcErr:  mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "file download fails due to user id is misssing",
			mockGrpcErr:  mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "file download succeeds with no file id found",
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
			expectedData: MMSpb.DownloadFileResponse{},
		},
		{
			name:         "file download fails due to internal error",
			mockGrpcErr:  mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToDownloadFile.Error(),
		},
		{
			name:         "file download succeeds",
			mockGrpcErr:  mocks.GrpcOpSuccess,
			expectedCode: codes.OK,
			expectedErr:  "",
			expectedData: MMSpb.DownloadFileResponse{
				FileName: "test-file.txt",
				MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
				Content:  []byte("test content"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{MockErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{}
			c := NewClient(mmsClient, svc)

			name, mimeType, content, err := c.DownloadFileGrpcHandler(context.Background(), "user-123", "file-id-123")

			status, _ := status.FromError(err)

			if tt.expectedCode != status.Code() {
				t.Errorf("expected %v got %v", tt.expectedCode, status.Code())
			}

			if tt.expectedCode == codes.OK {
				if name != tt.expectedData.FileName {
					t.Errorf("expected filename %v got %v", tt.expectedData.FileName, name)
				}
				if mimeType != tt.expectedData.MimeType {
					t.Errorf("expected mime type %v got %v", tt.expectedData.MimeType, mimeType)
				}
				if string(content) != string(tt.expectedData.Content) {
					t.Errorf("expected content %v got %v", tt.expectedData.Content, content)
				}
			}

			if tt.expectedErr != status.Message() {
				t.Errorf("expected error %v got %v", tt.expectedErr, status.Message())
			}
		})
	}
}
