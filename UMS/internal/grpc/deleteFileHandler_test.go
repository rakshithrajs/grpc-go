package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

func TestDeleteFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		fileID            string
		deleteDbErr       mocks.DbOperationError
		grpcErr           mocks.GrpcOperationError
		createRollbackErr mocks.DbOperationError
		expectedCode      int
		expectedErr       string
	}{
		{
			name:         "file deleted successfully",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "delete succeeds with no user file found",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			deleteDbErr:  mocks.DbOpNotFound,
			expectedCode: http.StatusOK,
		},
		{
			name:         "delete succeeds with no file found in MMS",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "delete fails when storage returns internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			deleteDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  storage.ErrFailedToDeleteUserFile.Error(),
		},
		{
			name:         "delete fails due to missing metadata",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "delete fails due to missing user id",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "delete fails due to internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  storage.ErrFailedToDeleteUserFile.Error(),
		},
		{
			name:         "delete fails due to internal error and rollback succeeds",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  storage.ErrFailedToDeleteUserFile.Error(),
		},
		{
			name:              "delete fails due to internal error and rollback fails",
			fileID:            "550e8400-e29b-41d4-a716-446655440000",
			grpcErr:           mocks.GrpcOpInternalError,
			createRollbackErr: mocks.DbOpRollbackFailure,
			expectedCode:      http.StatusInternalServerError,
			expectedErr:       handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{DeleteGrpcErr: tt.grpcErr}
			svc := &mocks.MockUserFilesService{DeleteUserFileError: tt.deleteDbErr, CreateUserFileError: tt.createRollbackErr}
			c := NewClient(mmsClient, svc)

			status, errMsg := c.DeleteFileGrpcHandler(context.Background(), "user-123", tt.fileID)

			mockUtils.CheckData(t, status, tt.expectedCode)
			mockUtils.CheckError(t, errMsg, tt.expectedErr)
		})
	}
}
