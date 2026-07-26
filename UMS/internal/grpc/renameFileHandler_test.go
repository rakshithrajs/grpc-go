package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
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
			name:         "file rename failed as file id is missing",
			fileID:       config.NullString,
			newName:      "renamed.txt",
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file rename failed as file id is whitespace",
			fileID:       "   ",
			newName:      "renamed.txt",
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file rename failed as new name is missing",
			fileID:       "file-id-123",
			newName:      config.NullString,
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrNewNameRequired.Error(),
		},
		{
			name:         "file rename failed as new name is whitespace",
			fileID:       "file-id-123",
			newName:      "   ",
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrNewNameRequired.Error(),
		},
		{
			name:         "file rename failed due to db internal error",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToUpdateUserFile.Error(),
		},
		{
			name:         "file rename succeeded with no old name found",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			updateDbErr:  mocks.DbOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "file rename failed due to missing metadata",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "file rename failed due to missing user id",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "file rename succeeded with no file id found",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "file rename fails as file name already exists",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpFileNameAlreadyExists,
			expectedCode: http.StatusConflict,
			expectedErr:  mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "file rename fails due to internal error",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToRenameFile.Error(),
		},
		{
			name:              "grpc rename internal error with rollback failure",
			fileID:            "file-id-123",
			newName:           "renamed.txt",
			grpcErr:           mocks.GrpcOpInternalError,
			updateRollbackErr: mocks.DbOpInternalError,
			expectedCode:      http.StatusInternalServerError,
			expectedErr:       handlerErrors.ErrFailedToRollback.Error(),
		},
		{
			name:         "rename succeeds",
			fileID:       "file-id-123",
			newName:      "renamed.txt",
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{MockErr: tt.grpcErr}
			svc := &mocks.MockUserFilesService{UpdateUserFileError: tt.updateDbErr, UpdateRollbackError: tt.updateRollbackErr}
			c := NewClient(mmsClient, svc)

			status, errMsg := c.RenameFileGrpcHandler(context.Background(), "user-123", tt.fileID, tt.newName)

			if tt.expectedCode != status {
				t.Errorf("expected %v got %v", tt.expectedCode, status)
			}

			if status != http.StatusOK {
				mockUtils.CheckData(t, errMsg, tt.expectedErr)
			}
		})
	}
}
