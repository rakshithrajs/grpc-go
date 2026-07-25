package grpc

import (
	"context"
	"net/http"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

func TestDownloadFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name         string
		fileID       string
		mockGrpcErr  mocks.GrpcOperationError
		expectedCode int
		expectedErr  string
		expectedData *MMSpb.DownloadFileResponse
	}{
		{
			name:         "file download fails as file id is missing",
			fileID:       config.NullString,
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file download fails as file id is whitespace",
			fileID:       "   ",
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file download fails as missing metadata",
			fileID:       "file-id-123",
			mockGrpcErr:  mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "file download fails due to user id is missing",
			fileID:       "file-id-123",
			mockGrpcErr:  mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "file download succeeds with no file id found",
			fileID:       "file-id-123",
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
			expectedData: &MMSpb.DownloadFileResponse{},
		},
		{
			name:         "file download fails due to internal error",
			fileID:       "file-id-123",
			mockGrpcErr:  mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerUtils.ErrFailedToDownloadFile.Error(),
		},
		{
			name:         "file download succeeds",
			fileID:       "file-id-123",
			mockGrpcErr:  mocks.GrpcOpSuccess,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
			expectedData: &MMSpb.DownloadFileResponse{
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

			resp, status, err := c.DownloadFileGrpcHandler(context.Background(), "user-123", tt.fileID)

			if tt.expectedCode != status {
				t.Errorf("expected %v got %v", tt.expectedCode, status)
			}

			if tt.expectedCode == http.StatusOK {
				mockUtils.CheckData(t, resp, tt.expectedData)
			}

			mockUtils.CheckData(t, err, tt.expectedErr)
		})
	}
}
