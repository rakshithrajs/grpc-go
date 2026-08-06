package grpcClient

import (
	"context"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestRenameFileGrpcClient(t *testing.T) {
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
			name:         "file renamed successfully",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "rename succeeds with no user file found",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "rename succeeds with no file found in MMS",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "rename fails when storage returns internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerErrors.ErrFailedToUpdateUserFile.Error(),
		},
		{
			name:         "rename fails due to missing metadata",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "rename fails due to missing user id",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "rename fails as file already exists",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpFileAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFileAlreadyExists.Error(),
		},
		{
			name:         "rename fails due to internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerErrors.ErrFailedToRenameFile.Error(),
		},
		{
			name:         "rename fails due to internal error and rollback succeeds",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerErrors.ErrFailedToRenameFile.Error(),
		},
		{
			name:              "rename fails due to internal error and rollback fails",
			fileID:            "550e8400-e29b-41d4-a716-446655440000",
			newName:           "renamed.txt",
			grpcErr:           mocks.GrpcOpInternalError,
			updateRollbackErr: mocks.DbOpInternalError,
			expectedCode:      codes.Internal,
			expectedErr:       handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{RenameGrpcErr: tt.grpcErr}
			svc := &mocks.MockUserFilesService{UpdateUserFileError: tt.updateDbErr, UpdateRollbackError: tt.updateRollbackErr}
			c := NewMMSClient(mmsClient, svc)

			err := c.RenameFileGrpcClient(context.Background(), "user-123", tt.fileID, tt.newName)

			status, _ := status.FromError(err)

			mockUtils.CheckData(t, status.Code(), tt.expectedCode)
			mockUtils.CheckError(t, status.Message(), tt.expectedErr)
		})
	}
}
