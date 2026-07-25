package handlers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

func TestRenameFileHandler(t *testing.T) {
	tests := []struct {
		name                string
		auth                bool
		fileID              string
		body                string
		UpdateUserFileError mocks.DbOperationError
		UpdateRollbackError mocks.DbOperationError
		mockGrpcErr         mocks.GrpcOperationError
		mockGrpcErr2        mocks.GrpcOperationError
		expectedCode        int
		expectedError       any
		expectedData        any
	}{
		{
			name:          "rename file fails due to missing auth",
			auth:          false,
			body:          `{"newName":"renamed.txt"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrUnauthorized.Error(),
		},
		{
			name:          "rename file fails due to missing fileID",
			auth:          true,
			fileID:        "",
			body:          `{"newName":"renamed.txt"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrFileIDRequired.Error(),
		},
		{
			name:          "rename file fails due to whitespace fileID",
			auth:          true,
			fileID:        "   ",
			body:          `{"newName":"renamed.txt"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrFileIDRequired.Error(),
		},
		{
			name:          "rename file fails due to invalid json",
			auth:          true,
			fileID:        "file-id-123",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrInvalidJSON.Error(),
		},
		{
			name:         "rename file fails due to empty newName",
			auth:         true,
			fileID:       "file-id-123",
			body:         `{"newName":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"newName": modelUtils.ErrNewNameRequired.Error(),
			},
		},
		{
			name:         "rename file fails due to whitespace newName",
			auth:         true,
			fileID:       "file-id-123",
			body:         `{"newName":"   "}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"newName": modelUtils.ErrNewNameRequired.Error(),
			},
		},
		{
			name:                "rename file fails due to db internal error",
			auth:                true,
			fileID:              "file-id-123",
			body:                `{"newName":"renamed.txt"}`,
			UpdateUserFileError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       handlerUtils.ErrFailedToRenameFile.Error(),
		},
		{
			name:                "rename file succeeds when file not found in db",
			auth:                true,
			fileID:              "file-id-123",
			body:                `{"newName":"renamed.txt"}`,
			UpdateUserFileError: mocks.DbOpNotFound,
			expectedCode:        http.StatusOK,
			expectedData:        map[string]string{"message": fileRenamedSuccessfully},
		},
		{
			name:          "rename file fails due to grpc internal error",
			auth:          true,
			fileID:        "file-id-123",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerUtils.ErrFailedToRenameFile.Error(),
		},
		{
			name:          "rename file fails due to grpc invalid argument",
			auth:          true,
			fileID:        "file-id-123",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpInvalidArgument,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrFileIDRequired.Error(),
		},
		{
			name:          "rename file fails due to grpc file name already exists",
			auth:          true,
			fileID:        "file-id-123",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpFileNameAlreadyExists,
			expectedCode:  http.StatusConflict,
			expectedError: mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "rename file succeeds when file not found in grpc",
			auth:         true,
			fileID:       "file-id-123",
			body:         `{"newName":"renamed.txt"}`,
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedData: map[string]string{"message": fileRenamedSuccessfully},
		},
		{
			name:                "rename file fails due to grpc rollback failure",
			auth:                true,
			fileID:              "file-id-123",
			body:                `{"newName":"renamed.txt"}`,
			mockGrpcErr:         mocks.GrpcOpInternalError,
			UpdateRollbackError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       handlerUtils.ErrFailedToRenameFile.Error(),
		},
		{
			name:         "rename file succeeds",
			auth:         true,
			fileID:       "file-id-123",
			body:         `{"newName":"renamed.txt"}`,
			expectedCode: http.StatusOK,
			expectedData: map[string]string{"message": fileRenamedSuccessfully},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTest(http.MethodPatch, "/api/files/:fileID/rename", tt.body, tt.auth)
			c.Params = []gin.Param{{Key: "fileID", Value: tt.fileID}}

			mmsClient := &mocks.MockMMSClient{MockErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{UpdateUserFileError: tt.UpdateUserFileError, UpdateRollbackError: tt.UpdateRollbackError}
			client := grpc.NewClient(mmsClient, svc)
			handler := NewUserFilesHandler(client, svc)

			handler.RenameFileHandler(c)

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
