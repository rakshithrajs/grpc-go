package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
)

func TestRenameFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		fileID            string
		newName           string
		grpcErr           mocks.GrpcOperationError
		updateDbErr       mocks.DbOperationError
		updateRollbackErr mocks.DbOperationError
		expectedCode      int
		expectedErr       string
	}{
		{
			name:         "file renamed successfully",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "rename succeeds with no user file found",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "rename succeeds with no file found in MMS",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "rename fails when storage returns internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToUpdateUserFile.Error(),
		},
		{
			name:         "rename fails due to missing metadata",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "rename fails due to missing user id",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "rename fails as file already exists",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpFileAlreadyExists,
			expectedCode: http.StatusConflict,
			expectedErr:  mocks.ErrFileAlreadyExists.Error(),
		},
		{
			name:         "rename fails due to internal error",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToRenameFile.Error(),
		},
		{
			name:         "rename fails due to internal error and rollback succeeds",
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToRenameFile.Error(),
		},
		{
			name:              "rename fails due to internal error and rollback fails",
			fileID:            "550e8400-e29b-41d4-a716-446655440000",
			newName:           "renamed.txt",
			grpcErr:           mocks.GrpcOpInternalError,
			updateRollbackErr: mocks.DbOpInternalError,
			expectedCode:      http.StatusInternalServerError,
			expectedErr:       handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{RenameGrpcErr: tt.grpcErr}
			svc := &mocks.MockUserFilesService{UpdateUserFileError: tt.updateDbErr, UpdateRollbackError: tt.updateRollbackErr}
			c := NewClient(mmsClient, svc)

			status, errMsg := c.RenameFileGrpcHandler(context.Background(), "user-123", tt.fileID, tt.newName)

			mockUtils.CheckData(t, status, tt.expectedCode)
			mockUtils.CheckError(t, errMsg, tt.expectedErr)
		})
	}
}
