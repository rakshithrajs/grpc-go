package grpcClient

import (
	"context"
	"testing"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDownloadFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name         string
		fileID       string
		mockGrpcErr  mocks.GrpcOperationError
		expectedCode codes.Code
		expectedErr  string
		expectedData *MMSpb.DownloadFileResponse
	}{
		{
			name:         "file downloaded successfully",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: codes.OK,
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
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
			expectedData: &MMSpb.DownloadFileResponse{},
		},
		{
			name:         "download fails due to missing metadata",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "download fails due to missing user id",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "download fails due to internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerErrors.ErrFailedToDownloadFile.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{DownloadGrpcErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{}
			c := NewClient(mmsClient, svc)

			resp, err := c.DownloadFileGrpcHandler(context.Background(), "user-123", tt.fileID)

			status, _ := status.FromError(err)

			mockUtils.CheckData(t, status.Code(), tt.expectedCode)
			mockUtils.CheckError(t, status.Message(), tt.expectedErr)

			if status.Code() == codes.OK {
				mockUtils.CheckData(t, resp, tt.expectedData)
			}
		})
	}
}
