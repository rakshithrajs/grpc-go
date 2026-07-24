package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func TestUploadFileHandler(t *testing.T) {
	tests := []struct {
		name                string
		auth                bool
		withFile            bool
		CreateUserFileError mocks.DbOperationError
		uploadGrpcErr       mocks.GrpcOperationError
		deleteGrpcErr       mocks.GrpcOperationError
		expectedCode        int
		expectedError       any
		expectedData        any
	}{
		{
			name:          "upload file fails due to missing auth",
			auth:          false,
			withFile:      true,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "upload file fails due to missing file",
			auth:          true,
			withFile:      false,
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrFileIsRequired.Error(),
		},
		{
			name:          "upload file fails due to grpc internal error",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrFailedToUploadFile.Error(),
		},
		{
			name:          "upload file fails due to grpc file name already exists",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpFileNameAlreadyExists,
			expectedCode:  http.StatusConflict,
			expectedError: mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:          "upload file fails due to grpc file path already exists",
			auth:          true,
			withFile:      true,
			uploadGrpcErr: mocks.GrpcOpFilePathAlreadyExists,
			expectedCode:  http.StatusConflict,
			expectedError: mocks.ErrFilePathAlreadyExists.Error(),
		},
		{
			name:                "upload file fails due to db duplicate file",
			auth:                true,
			withFile:            true,
			CreateUserFileError: mocks.DbOpDuplicateFile,
			expectedCode:        http.StatusConflict,
			expectedError:       utils.ErrUserFileAlreadyExists.Error(),
		},
		{
			name:                "upload file fails due to db internal error",
			auth:                true,
			withFile:            true,
			CreateUserFileError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       utils.ErrFailedToUploadFile.Error(),
		},
		{
			name:                "upload file fails due to grpc rollback failure",
			auth:                true,
			withFile:            true,
			CreateUserFileError: mocks.DbOpInternalError,
			deleteGrpcErr:       mocks.GrpcOpRollbackFailure,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       utils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload file succeeds",
			auth:         true,
			withFile:     true,
			expectedCode: http.StatusCreated,
			expectedData: map[string]any{
				"file": map[string]any{
					"id":       "file-id-123",
					"fileName": "test.txt",
					"fileSize": float64(12),
					"mimeType": "text/plain",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c *gin.Context
			var w *httptest.ResponseRecorder
			if tt.withFile {
				c, w = mocks.SetUpGinTestMultipart("test content", tt.auth)
			} else {
				c, w = mocks.SetUpGinTest(http.MethodPost, "/api/files/upload", "", tt.auth)
			}

			mmsClient := &mocks.MockMMSClient{MockErr: tt.uploadGrpcErr, MockDeleteErr: tt.deleteGrpcErr}
			svc := &mocks.MockUserFilesService{CreateUserFileError: tt.CreateUserFileError}
			client := grpc.NewClient(mmsClient, svc)
			handler := NewUserFilesHandler(client, svc)

			handler.UploadFileHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedError != nil {
				mocks.CheckError(t, w, tt.expectedError)
			}

			if tt.expectedData != nil {
				mocks.CheckData(t, w, tt.expectedData)
			}
		})
	}
}
