package grpcClient

import (
	"context"
	"testing"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUploadFileGrpcClient(t *testing.T) {
	tests := []struct {
		name              string
		fileName          string
		content           []byte
		grpcErr           mocks.GrpcOperationError
		deleteGrpcErr     mocks.GrpcOperationError
		createDbErr       mocks.DbOperationError
		uploadReturnEmpty bool
		expectedCode      codes.Code
		expectedErr       string
		expectedFile      *models.File
	}{
		{
			name:         "file uploaded successfully",
			fileName:     "test.txt",
			content:      []byte("test content"),
			expectedCode: codes.OK,
			expectedFile: &models.File{
				ID:       "550e8400-e29b-41d4-a716-446655440000",
				FileName: "test.txt",
				FileSize: 12,
				MimeType: "text/plain",
			},
		},
		{
			name:         "upload fails due to missing metadata",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "upload fails due to missing user id",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "upload fails due to internal error",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerErrors.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails as file already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFileAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFileAlreadyExists.Error(),
		},
		{
			name:         "upload fails when user file mapping already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpDuplicateFile,
			expectedCode: codes.AlreadyExists,
			expectedErr:  handlerErrors.ErrUserFileAlreadyExists.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback succeeds",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerErrors.ErrFailedToUploadFile.Error(),
		},
		{
			name:          "upload fails due to db internal error and rollback fails",
			fileName:      "test.txt",
			content:       []byte("test content"),
			createDbErr:   mocks.DbOpInternalError,
			deleteGrpcErr: mocks.GrpcOpRollbackFailure,
			expectedCode:  codes.Internal,
			expectedErr:   handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{UploadGrpcErr: tt.grpcErr, DeleteGrpcErr: tt.deleteGrpcErr}
			svc := &mocks.MockUserFilesService{CreateUserFileError: tt.createDbErr}
			c := NewMMSClient(mmsClient, svc)

			file, err := c.UploadFileGrpcClient(context.Background(), "user-123", tt.fileName, tt.content)

			status, _ := status.FromError(err)

			mockUtils.CheckData(t, status.Code(), tt.expectedCode)
			mockUtils.CheckError(t, status.Message(), tt.expectedErr)

			if status.Code() == codes.OK {
				mockUtils.CheckData(t, file, tt.expectedFile)
			}
		})
	}
}
