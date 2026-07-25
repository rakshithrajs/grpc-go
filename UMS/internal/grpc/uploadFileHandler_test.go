package grpc

import (
	"context"
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
)

func TestClient_UploadFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		fileName          string
		content           []byte
		grpcErr           mocks.GrpcOperationError
		createDbErr       mocks.DbOperationError
		uploadReturnEmpty bool
		expectedCode      int
		expectedErr       string
		expectedFile      *models.File
	}{
		{
			name:         "upload fails when file name is empty",
			content:      []byte("test content"),
			fileName:     config.NullString,
			expectedCode: http.StatusBadRequest,
			expectedErr:  handlerUtils.ErrFileNameRequired.Error(),
		},
		{
			name:         "upload fails when file name is whitespace",
			content:      []byte("test content"),
			fileName:     "   ",
			expectedCode: http.StatusBadRequest,
			expectedErr:  handlerUtils.ErrFileNameRequired.Error(),
		},
		{
			name:         "upload fails when content is empty",
			fileName:     "test.txt",
			content:      []byte{},
			expectedCode: http.StatusBadRequest,
			expectedErr:  handlerUtils.ErrFileIsRequired.Error(),
		},
		{
			name:         "upload fails due to missing metadata",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "upload fails due to missing user id",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "upload fails due to internal error",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerUtils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails as file name already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFileNameAlreadyExists,
			expectedCode: http.StatusConflict,
			expectedErr:  mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "upload fails as file path already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFilePathAlreadyExists,
			expectedCode: http.StatusConflict,
			expectedErr:  mocks.ErrFilePathAlreadyExists.Error(),
		},
		{
			name:              "upload fails when grpc returns empty file",
			fileName:          "test.txt",
			content:           []byte("test content"),
			uploadReturnEmpty: true,
			expectedCode:      http.StatusInternalServerError,
			expectedErr:       handlerUtils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails when user file mapping already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpDuplicateFile,
			expectedCode: http.StatusConflict,
			expectedErr:  handlerUtils.ErrUserFileAlreadyExists.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback succeeds",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerUtils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback fails",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpRollbackFailure,
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerUtils.ErrFailedToRollback.Error(),
		},
		{
			name:         "upload succeeds",
			fileName:     "test.txt",
			content:      []byte("test content"),
			expectedCode: http.StatusOK,
			expectedFile: &models.File{
				ID:       "file-id-123",
				FileName: "test.txt",
				FileSize: 12,
				MimeType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{
				MockErr:           tt.grpcErr,
				UploadReturnEmpty: tt.uploadReturnEmpty,
			}
			svc := &mocks.MockUserFilesService{CreateUserFileError: tt.createDbErr}
			c := NewClient(mmsClient, svc)

			file, status, errMsg := c.UploadFileGrpcHandler(context.Background(), "user-123", tt.fileName, tt.content)

			if tt.expectedCode != status {
				t.Errorf("expected code %v got %v", tt.expectedCode, status)
			}

			if status != http.StatusOK {
				if tt.expectedErr != errMsg {
					t.Errorf("expected error %v got %v", tt.expectedErr, errMsg)
				}
				return
			}

			mockUtils.CheckData(t, file, tt.expectedFile)
		})
	}
}
