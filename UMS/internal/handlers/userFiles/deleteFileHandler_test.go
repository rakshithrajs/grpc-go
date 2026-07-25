package handlers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func TestDeleteFileHandler(t *testing.T) {
	tests := []struct {
		name                string
		auth                bool
		fileID              string
		DeleteUserFileError mocks.DbOperationError
		CreateUserFileError mocks.DbOperationError
		mockGrpcErr         mocks.GrpcOperationError
		expectedCode        int
		expectedError       any
		expectedData        any
	}{
		{
			name:          "delete file fails due to missing auth",
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "delete file fails due to missing fileID",
			auth:          true,
			fileID:        "",
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrFileIDRequired.Error(),
		},
		{
			name:          "delete file fails due to whitespace fileID",
			auth:          true,
			fileID:        "   ",
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrFileIDRequired.Error(),
		},
		{
			name:                "delete file fails due to db internal error",
			auth:                true,
			fileID:              "file-id-123",
			DeleteUserFileError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       utils.ErrFailedToDeleteFile.Error(),
		},
		{
			name:                "delete file succeeds when file not found in db",
			auth:                true,
			fileID:              "file-id-123",
			DeleteUserFileError: mocks.DbOpNotFound,
			expectedCode:        http.StatusOK,
			expectedData:        map[string]string{"message": deleteFileSuccessMsg},
		},
		{
			name:          "delete file fails due to grpc internal error",
			auth:          true,
			fileID:        "file-id-123",
			mockGrpcErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrFailedToDeleteFile.Error(),
		},
		{
			name:          "delete file fails due to grpc invalid argument",
			auth:          true,
			fileID:        "file-id-123",
			mockGrpcErr:   mocks.GrpcOpInvalidArgument,
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrFileIDRequired.Error(),
		},
		{
			name:         "delete file succeeds when file not found in grpc",
			auth:         true,
			fileID:       "file-id-123",
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedData: map[string]string{"message": deleteFileSuccessMsg},
		},
		{
			name:                "delete file fails due to grpc rollback failure",
			auth:                true,
			fileID:              "file-id-123",
			mockGrpcErr:         mocks.GrpcOpInternalError,
			CreateUserFileError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       utils.ErrFailedToDeleteFile.Error(),
		},
		{
			name:         "delete file succeeds",
			auth:         true,
			fileID:       "file-id-123",
			expectedCode: http.StatusOK,
			expectedData: map[string]string{"message": deleteFileSuccessMsg},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodDelete, "/api/files/:fileID", "", tt.auth)
			c.Params = []gin.Param{{Key: "fileID", Value: tt.fileID}}

			mmsClient := &mocks.MockMMSClient{MockErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{DeleteUserFileError: tt.DeleteUserFileError, CreateUserFileError: tt.CreateUserFileError}
			client := grpc.NewClient(mmsClient, svc)
			handler := NewUserFilesHandler(client, svc)

			handler.DeleteFileHandler(c)

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
