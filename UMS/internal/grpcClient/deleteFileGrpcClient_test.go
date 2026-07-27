package grpcClient

import (
	"context"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestDeleteFileGrpcClient(t *testing.T) {
	tests := []struct {
		name              string
		fileID            string
		deleteDbErr       mocks.DbOperationError
		grpcErr           mocks.GrpcOperationError
		createRollbackErr mocks.DbOperationError
		expectedCode      codes.Code
		expectedErr       string
	}{
		{
			name:         "file deleted successfully",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "delete succeeds with no user file found",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			deleteDbErr:  mocks.DbOpNotFound,
			expectedCode: codes.OK,
		},
		{
			name:         "delete succeeds with no file found in MMS",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpNotFound,
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
		{
			name:         "delete fails when storage returns internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			deleteDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  storage.ErrFailedToDeleteUserFile.Error(),
		},
		{
			name:         "delete fails due to missing metadata",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "delete fails due to missing user id",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "delete fails due to internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToDeleteFile.Error(),
		},
		{
			name:         "delete fails due to internal error and rollback succeeds",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToDeleteFile.Error(),
		},
		{
			name:              "delete fails due to internal error and rollback fails",
			fileID:            "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:           mocks.GrpcOpInternalError,
			createRollbackErr: mocks.DbOpRollbackFailure,
			expectedCode:      codes.Internal,
			expectedErr:       handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{DeleteGrpcErr: tt.grpcErr}
			svc := &mocks.MockUserFilesService{DeleteUserFileError: tt.deleteDbErr, CreateUserFileError: tt.createRollbackErr}
			c := NewClient(mmsClient, svc)

			err := c.DeleteFileGrpcClient(context.Background(), "user-123", tt.fileID)

			status, _ := status.FromError(err)

			mockUtils.CheckData(t, status.Code(), tt.expectedCode)
			mockUtils.CheckError(t, status.Message(), tt.expectedErr)
		})
	}
}
