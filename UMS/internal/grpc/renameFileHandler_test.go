package grpc

import (
	"context"
	"testing"

	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRenameFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		fileID            string
		newName           string
		grpcErr           mocks.GrpcOperationError
		updateDbErr       mocks.DbOperationError
		updateRollbackErr mocks.DbOperationError
		expectedCode      codes.Code
		expectedErr       string
	}{
		{
			name:         "file rename failed as file id is missing",
			fileID:       "",
			newName:      "renamed.txt",
			expectedCode: codes.InvalidArgument,
			expectedErr:  handlerUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file rename failed as file id is whitespace",
			fileID:       "   ",
			newName:      "renamed.txt",
			expectedCode: codes.InvalidArgument,
			expectedErr:  handlerUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file rename failed as new name is missing",
			fileID:       "file-id-123",
			newName:      "",
			expectedCode: codes.InvalidArgument,
			expectedErr:  modelUtils.ErrNewNameRequired.Error(),
		},
		{
			name:         "file rename failed as new name is whitespace",
			fileID:       "file-id-123",
			newName:      "   ",
			expectedCode: codes.InvalidArgument,
			expectedErr:  modelUtils.ErrNewNameRequired.Error(),
		},
		{
			name:         "file rename failed due to db internal error",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerUtils.ErrFailedToUpdateUserFile.Error(),
		},
		{
			name:         "file rename succeeded with no old name found",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:         "file rename failed due to missing metadata",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "file rename failed due to missing user id",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "file rename succesded with no file id found",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  "",
		},
		{
			name:         "file rename fails as file name already exists",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpFileNameAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "file rename fails due to internal error",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToRenameFile.Error(),
		},
		{
			name:              "grpc rename internal error with rollback failure",
			fileID:            "file-id-123",
			newName:           "renamed.txt",
			grpcErr:           mocks.GrpcOpInternalError,
			updateRollbackErr: mocks.DbOpInternalError,
			expectedCode:      codes.Internal,
			expectedErr:       handlerUtils.ErrFailedToRollback.Error(),
		},
		{
			name:         "rename succeeds",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			expectedCode: codes.OK,
			expectedErr:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{MockErr: tt.grpcErr}
			svc := &mocks.MockUserFilesService{UpdateUserFileError: tt.updateDbErr, UpdateRollbackError: tt.updateRollbackErr}
			c := NewClient(mmsClient, svc)

			err := c.RenameFileGrpcHandler(context.Background(), "user-123", tt.fileID, tt.newName)

			status, _ := status.FromError(err)

			if tt.expectedCode != status.Code() {
				t.Errorf("expected %v got %v", tt.expectedCode, status.Code())
			}

			if status.Code() != codes.OK {
				mockUtils.CheckData(t, status.Message(), tt.expectedErr)
			}
		})
	}
}
