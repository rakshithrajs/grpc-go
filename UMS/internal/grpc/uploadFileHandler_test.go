package grpc

import (
	"context"
	"net/http"
	"testing"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
)

func TestUploadFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		fileName          string
		content           []byte
		grpcErr           mocks.GrpcOperationError
		deleteGrpcErr     mocks.GrpcOperationError
		createDbErr       mocks.DbOperationError
		uploadReturnEmpty bool
		expectedCode      int
		expectedErr       string
		expectedFile      *models.File
	}{
		{
			name:         "file uploaded successfully",
			fileName:     "test.txt",
			content:      []byte("test content"),
			expectedCode: http.StatusCreated,
			expectedFile: &models.File{
				ID:       "550e8400-e29b-41d4-a716-446655440000",
				FileName: "test.txt",
				FileSize: 12,
				MimeType: "text/plain",
			},
		},
		{
			name:         "upload fails due to missing metadata",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "upload fails due to missing user id",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: http.StatusUnauthorized,
			expectedErr:  handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "upload fails due to internal error",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails as file already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFileAlreadyExists,
			expectedCode: http.StatusConflict,
			expectedErr:  mocks.ErrFileAlreadyExists.Error(),
		},
		{
			name:         "upload fails when user file mapping already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpDuplicateFile,
			expectedCode: http.StatusConflict,
			expectedErr:  handlerErrors.ErrUserFileAlreadyExists.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback succeeds",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: http.StatusInternalServerError,
			expectedErr:  handlerErrors.ErrFailedToUploadFile.Error(),
		},
		{
			name:          "upload fails due to db internal error and rollback fails",
			fileName:      "test.txt",
			content:       []byte("test content"),
			createDbErr:   mocks.DbOpInternalError,
			deleteGrpcErr: mocks.GrpcOpRollbackFailure,
			expectedCode:  http.StatusInternalServerError,
			expectedErr:   handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{UploadGrpcErr: tt.grpcErr, DeleteGrpcErr: tt.deleteGrpcErr}
			svc := &mocks.MockUserFilesService{CreateUserFileError: tt.createDbErr}
			c := NewClient(mmsClient, svc)

			file, status, errMsg := c.UploadFileGrpcHandler(context.Background(), "user-123", tt.fileName, tt.content)

			mockUtils.CheckData(t, status, tt.expectedCode)
			mockUtils.CheckError(t, errMsg, tt.expectedErr)

			if status == http.StatusCreated {
				mockUtils.CheckData(t, file, tt.expectedFile)
			}
		})
	}
}
