package grpc

import (
	"context"
	"net/http"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
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
			name:         "file downloaded successfully",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
			expectedData: &MMSpb.DownloadFileResponse{
				FileName: "test-file.txt",
				MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
				Content:  []byte("test content"),
			},
		},
		{
			name:         "download succeeds with no file found",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
			expectedData: &MMSpb.DownloadFileResponse{},
		},
		{
			name:         "download fails due to missing metadata",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "download fails due to missing user id",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "download fails due to internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToDownloadFile.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{DownloadGrpcErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{}
			c := NewClient(mmsClient, svc)

			resp, status, err := c.DownloadFileGrpcHandler(context.Background(), "user-123", tt.fileID)

			mockUtils.CheckData(t, status, tt.expectedCode)
			mockUtils.CheckError(t, err, tt.expectedErr)

			if status == http.StatusOK {
				mockUtils.CheckData(t, resp, tt.expectedData)
			}
		})
	}
}
