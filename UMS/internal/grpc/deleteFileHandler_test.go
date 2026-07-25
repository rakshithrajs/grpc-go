package grpc

import (
	"context"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name         string
		fileID       string
		deleteDbErr  mocks.DbOperationError
		createDbErr  mocks.DbOperationError
		GrpcErr      mocks.GrpcOperationError
		expectedCode codes.Code
		expectedErr  string
	}{
		{
			name:         "file deletion failed as file id is missing",
			fileID:       config.NullString,
			expectedCode: codes.InvalidArgument,
			expectedErr:  handlerUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file deletion failed as file id is whitespace",
			fileID:       "   ",
			expectedCode: codes.InvalidArgument,
			expectedErr:  handlerUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file deletion failed due to db internal error",
			fileID:       "file-id-123",
			deleteDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerUtils.ErrFailedToDeleteUserFile.Error(),
		},
		{
			name:         "file deletion succeeds but file not found in db",
			fileID:       "file-id-123",
			deleteDbErr:  mocks.DbOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "file deleted failed due to missing metadata",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "file deleted failed due to missing user id",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "file deletion failed due to grpc internal error",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerUtils.ErrFailedToDeleteFile.Error(),
		},
		{
			name:         "file deletion succeeds but file not found in grpc",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "file deletion failed due to grpc internal error with rollback failure",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpInternalError,
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerUtils.ErrFailedToRollback.Error(),
		},
		{
			name:         "file deletion succeeds",
			fileID:       "file-id-123",
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{MockErr: tt.GrpcErr}
			svc := &mocks.MockUserFilesService{DeleteUserFileError: tt.deleteDbErr, CreateUserFileError: tt.createDbErr}
			c := NewClient(mmsClient, svc)

			err := c.DeleteFileGrpcHandler(context.Background(), "user-123", tt.fileID)

			status, _ := status.FromError(err)

			if status.Code() != tt.expectedCode {
				t.Errorf("expected code = %v, got %v", tt.expectedCode, status.Code())
			}

			if status.Code() != codes.OK {
				mockUtils.CheckData(t, status.Message(), tt.expectedErr)
			}
		})
	}
}
