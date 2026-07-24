package grpc

import (
	"context"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name         string
		deleteDbErr  mocks.DbOperationError
		createDbErr  mocks.DbOperationError
		GrpcErr      mocks.GrpcOperationError
		expectedCode codes.Code
		expectedErr  string
	}{
		{
			name:         "file deletion failed due to db internal error",
			deleteDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  utils.ErrFailedToDeleteUserFile.Error(),
		},
		{
			name:         "file deletion succeeds but file not found in db",
			deleteDbErr:  mocks.DbOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:         "file deletion failed due to grpc internal error",
			GrpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  utils.ErrFailedToDeleteFile.Error(),
		},
		{
			name:         "file deletion succeeds but file not found in grpc",
			GrpcErr:      mocks.GrpcOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:         "file deletion failed due to grpc internal error with rollback failure",
			GrpcErr:      mocks.GrpcOpInternalError,
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  utils.ErrFailedToRollback.Error(),
		},
		{
			name:         "file deletion succeeds",
			expectedCode: codes.OK,
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{MockErr: tt.GrpcErr}
			svc := &mocks.MockUserFilesService{DeleteUserFileError: tt.deleteDbErr, CreateUserFileError: tt.createDbErr}
			c := NewClient(mmsClient, svc)

			err := c.DeleteFileGrpcHandler(context.Background(), "user-123", "file-id-123")

			status, _ := status.FromError(err)

			if status.Code() != tt.expectedCode {
				t.Errorf("expected code = %v, got %v", tt.expectedCode, status.Code())
			}

			if status.Code() != codes.OK {
				if status.Message() != tt.expectedErr {
					t.Errorf("expected error = %v, got %v", tt.expectedErr, status.Message())
				}
			}
		})
	}
}
