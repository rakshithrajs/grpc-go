package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
)

func TestUploadFileHandler(t *testing.T) {
	tests := []struct {
		name                string
		auth                bool
		withFile            bool
		fileName            string
		fileContent         string
		emptyContent        bool
		CreateUserFileError mocks.DbOperationError
		uploadGrpcErr       mocks.GrpcOperationError
		deleteGrpcErr       mocks.GrpcOperationError
		expectedCode        int
		expectedError       any
		expectedData        any
	}{
		{
			name:         "file uploaded successfully",
			auth:         true,
			withFile:     true,
			expectedCode: http.StatusCreated,
			expectedData: map[string]any{
				"file": map[string]any{
					"id":       "550e8400-e29b-41d4-a716-446655440000",
					"fileName": "test.txt",
					"fileSize": float64(12),
					"mimeType": "text/plain",
				},
			},
		},
		{
			name:          "upload file fails due to missing auth",
			auth:          false,
			withFile:      true,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "upload file fails due to missing file",
			auth:          true,
			withFile:      false,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerErrors.ErrFileIsRequired.Error(),
		},
		{
			name:          "upload file fails due to empty file content",
			auth:          true,
			withFile:      true,
			emptyContent:  true,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerErrors.ErrEmptyFileContent.Error(),
		},
		{
			name:          "upload file fails due to whitespace file name",
			auth:          true,
			withFile:      true,
			fileName:      "   ",
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerErrors.ErrFileNameRequired.Error(),
		},
		{
			name:          "upload file fails due to grpc missing metadata",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpMissingMetadata,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "upload file fails due to grpc missing userID",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpMissingUserID,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "upload file fails due to grpc internal error",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerErrors.ErrFailedToUploadFile.Error(),
		},
		{
			name:          "upload file fails due to grpc file already exists",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpFileAlreadyExists,
			expectedCode:  http.StatusConflict,
			expectedError: mocks.ErrFileAlreadyExists.Error(),
		},
		{
			name:                "upload file fails due to db duplicate file",
			auth:                true,
			withFile:            true,
			CreateUserFileError: mocks.DbOpDuplicateFile,
			expectedCode:        http.StatusConflict,
			expectedError:       handlerErrors.ErrUserFileAlreadyExists.Error(),
		},
		{
			name:                "upload file fails due to db internal error",
			auth:                true,
			withFile:            true,
			CreateUserFileError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       handlerErrors.ErrFailedToUploadFile.Error(),
		},
		{
			name:          "upload file fails due to grpc internal and rollback success",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerErrors.ErrFailedToUploadFile.Error(),
		},
		{
			name:                "upload file fails due to grpc rollback failure",
			auth:                true,
			withFile:            true,
			CreateUserFileError: mocks.DbOpInternalError,
			deleteGrpcErr:       mocks.GrpcOpRollbackFailure,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *gin.Context
			var w *httptest.ResponseRecorder
			if tt.withFile {
				content := tt.fileContent
				if content == config.NullString && !tt.emptyContent {
					content = "test content"
				}
				name := tt.fileName
				if name == config.NullString {
					name = "test.txt"
				}
				c, w = mockUtils.SetUpGinTestMultipart(content, name, tt.auth)
			} else {
				c, w = mockUtils.SetUpGinTest(http.MethodPost, "/api/files/upload", config.NullString, tt.auth)
			}

			mmsClient := &mocks.MockMMSClient{UploadGrpcErr: tt.uploadGrpcErr, DeleteGrpcErr: tt.deleteGrpcErr}
			svc := &mocks.MockUserFilesService{CreateUserFileError: tt.CreateUserFileError}
			client := grpc.NewClient(mmsClient, svc)
			handler := NewUserFilesHandler(client, svc)

			handler.UploadFileHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedError != nil {
				mockUtils.CheckError(t, w, tt.expectedError)
			}

			if tt.expectedData != nil {
				mockUtils.CheckData(t, w, tt.expectedData)
			}
		})
	}
}
