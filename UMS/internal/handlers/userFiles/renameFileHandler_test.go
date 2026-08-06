package handlers

import (
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpc "github.com/rakshithrajs/cloud/UMS/internal/grpcClient"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
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
		expectedCode        int
		expectedError       any
		expectedData        any
	}{
		{
			name:         "file renamed successfully",
			auth:         true,
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			body:         `{"newName":"renamed.txt"}`,
			expectedCode: http.StatusOK,
			expectedData: map[string]string{"message": fileRenamedSuccessfully},
		},
		{
			name:          "rename file fails due to missing auth",
			auth:          false,
			body:          `{"newName":"renamed.txt"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "rename file fails due to invalid json",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerErrors.ErrInvalidJSON.Error(),
		},
		{
			name:         "rename file fails due to missing fileID",
			auth:         true,
			fileID:       config.NullString,
			body:         `{"newName":"renamed.txt"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"fileID": modelUtils.ErrFileIDRequired.Error(),
			},
		},
		{
			name:         "rename file fails due to whitespace fileID",
			auth:         true,
			fileID:       "   ",
			body:         `{"newName":"renamed.txt"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"fileID": modelUtils.ErrFileIDRequired.Error(),
			},
		},
		{
			name:         "rename file fails due to invalid fileID",
			auth:         true,
			fileID:       "invalid-uuid",
			body:         `{"newName":"renamed.txt"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"fileID": modelUtils.ErrFileIDInvalidUUID.Error(),
			},
		},
		{
			name:         "rename file fails due to empty newName",
			auth:         true,
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			body:         `{"newName":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"newName": modelUtils.ErrNewNameRequired.Error(),
			},
		},
		{
			name:         "rename file fails due to whitespace newName",
			auth:         true,
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			body:         `{"newName":"   "}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"newName": modelUtils.ErrNewNameRequired.Error(),
			},
		},
		{
			name:         "rename file fails due to newName exceeding max length",
			auth:         true,
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			body:         `{"newName":"` + strings.Repeat("a", 256) + `"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"newName": modelUtils.ErrNameTooLong.Error(),
			},
		},
		{
			name:                "rename file fails due to db internal error",
			auth:                true,
			fileID:              "550e8400-e29b-41d4-a716-446655440000",
			body:                `{"newName":"renamed.txt"}`,
			UpdateUserFileError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       handlerErrors.ErrFailedToUpdateUserFile.Error(),
		},
		{
			name:                "rename file succeeds when file not found in db",
			auth:                true,
			fileID:              "550e8400-e29b-41d4-a716-446655440000",
			body:                `{"newName":"renamed.txt"}`,
			UpdateUserFileError: mocks.DbOpNotFound,
			expectedCode:        http.StatusOK,
			expectedData:        map[string]string{"message": fileRenamedSuccessfully},
		},
		{
			name:          "rename file fails due to grpc missing metadata",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpMissingMetadata,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "rename file fails due to grpc missing userID",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpMissingUserID,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "rename file fails due to grpc internal error",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerErrors.ErrFailedToRenameFile.Error(),
		},
		{
			name:          "rename file fails due to grpc file already exists",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpFileAlreadyExists,
			expectedCode:  http.StatusConflict,
			expectedError: mocks.ErrFileAlreadyExists.Error(),
		},
		{
			name:         "rename file succeeds when file not found in grpc",
			auth:         true,
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			body:         `{"newName":"renamed.txt"}`,
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedData: map[string]string{"message": fileRenamedSuccessfully},
		},
		{
			name:          "rename file fails due to grpc internal error and rollback succeeds",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			body:          `{"newName":"renamed.txt"}`,
			mockGrpcErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerErrors.ErrFailedToRenameFile.Error(),
		},
		{
			name:                "rename file fails due to grpc rollback failure",
			auth:                true,
			fileID:              "550e8400-e29b-41d4-a716-446655440000",
			body:                `{"newName":"renamed.txt"}`,
			mockGrpcErr:         mocks.GrpcOpInternalError,
			UpdateRollbackError: mocks.DbOpInternalError,
			expectedCode:        http.StatusInternalServerError,
			expectedError:       handlerErrors.ErrFailedToRollback.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTest(http.MethodPatch, "/api/files/:fileid/rename", tt.body, tt.auth)
			c.Params = []gin.Param{{Key: "fileid", Value: tt.fileID}}

			mmsClient := &mocks.MockMMSClient{RenameGrpcErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{UpdateUserFileError: tt.UpdateUserFileError, UpdateRollbackError: tt.UpdateRollbackError}
			client := grpc.NewMMSClient(mmsClient, svc)
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
