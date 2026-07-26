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

func TestDeleteFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name         string
		fileID       string
		deleteDbErr  mocks.DbOperationError
		createDbErr  mocks.DbOperationError
		GrpcErr      mocks.GrpcOperationError
		expectedCode int
		expectedErr  string
	}{
		{
			name:         "file deletion failed as file id is missing",
			fileID:       config.NullString,
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file deletion failed as file id is whitespace",
			fileID:       "   ",
			expectedCode: http.StatusBadRequest,
			expectedErr:  modelUtils.ErrFileIDRequired.Error(),
		},
		{
			name:         "file deletion failed due to db internal error",
			fileID:       "file-id-123",
			deleteDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToDeleteUserFile.Error(),
		},
		{
			name:         "file deletion succeeds but file not found in db",
			fileID:       "file-id-123",
			deleteDbErr:  mocks.DbOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "file deletion failed due to missing metadata",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "file deletion failed due to missing user id",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "file deletion failed due to grpc internal error",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToDeleteFile.Error(),
		},
		{
			name:         "file deletion succeeds but file not found in grpc",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
		{
			name:         "file deletion failed due to grpc internal error with rollback failure",
			fileID:       "file-id-123",
			GrpcErr:      mocks.GrpcOpInternalError,
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToRollback.Error(),
		},
		{
			name:         "file deletion succeeds",
			fileID:       "file-id-123",
			expectedCode: http.StatusOK,
			expectedErr:  config.NullString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{MockErr: tt.GrpcErr}
			svc := &mocks.MockUserFilesService{DeleteUserFileError: tt.deleteDbErr, CreateUserFileError: tt.createDbErr}
			c := NewClient(mmsClient, svc)

			status, err := c.DeleteFileGrpcHandler(context.Background(), "user-123", tt.fileID)

			if status != tt.expectedCode {
				t.Errorf("expected code = %v, got %v", tt.expectedCode, status)
			}

			if status != http.StatusOK {
				mockUtils.CheckData(t, err, tt.expectedErr)
			}
		})
	}
}
