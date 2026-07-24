package grpc

import (
	"context"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRenameFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		grpcErr           mocks.GrpcOperationError
		updateDbErr       mocks.DbOperationError
		updateRollbackErr mocks.DbOperationError
		expectedCode      codes.Code
		expectedErr       string
	}{
		{
			name:         "file rename failed due to db internal error",
			updateDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  utils.ErrFailedToUpdateUserFile.Error(),
		},
		{
			name:         "file rename succeeded with no old name found",
			updateDbErr:  mocks.DbOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:         "file rename failed due to missing metadata",
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "file rename failed due to missing user id",
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "file rename succesded with no file id found",
			grpcErr:      mocks.GrpcOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:         "file rename fails as file name already exists",
			grpcErr:      mocks.GrpcOpFileNameAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "file rename fails due to internal error",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToRenameFile.Error(),
		},
		{
			name:              "grpc rename internal error with rollback failure",
			grpcErr:           mocks.GrpcOpInternalError,
			updateRollbackErr: mocks.DbOpInternalError,
			expectedCode:      codes.Internal,
			expectedErr:       utils.ErrFailedToRollback.Error(),
		},
		{
			name:         "rename succeeds",
			expectedCode: codes.OK,
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{MockErr: tt.grpcErr}
			svc := &mocks.MockUserFilesService{UpdateUserFileError: tt.updateDbErr, UpdateRollbackError: tt.updateRollbackErr}
			c := NewClient(mmsClient, svc)

			err := c.RenameFileGrpcHandler(context.Background(), "user-123", "file-id-123", "renamed.txt")

			status, _ := status.FromError(err)

			if tt.expectedCode != status.Code() {
				t.Errorf("expected %v got %v", tt.expectedCode, status.Code())
			}

			if status.Code() != codes.OK {
				if tt.expectedErr != status.Message() {
					t.Errorf("expected %v got %v", tt.expectedErr, status.Message())
				}
			}
		})
	}
}
