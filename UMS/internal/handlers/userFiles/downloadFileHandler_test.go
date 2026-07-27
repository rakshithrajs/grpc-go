package handlers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

func TestDownloadFileHandler(t *testing.T) {
	tests := []struct {
		name          string
		auth          bool
		fileID        string
		mockGrpcErr   mocks.GrpcOperationError
		expectedCode  int
		expectedError any
		expectedBody  string
		expectedType  string
	}{
		{
			name:         "file downloaded successfully",
			auth:         true,
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			expectedCode: http.StatusOK,
			expectedBody: "test content",
			expectedType: "text/plain",
		},
		{
			name:          "download file fails due to missing auth",
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:         "download file fails due to missing fileID",
			auth:         true,
			fileID:       config.NullString,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"fileID": modelUtils.ErrFileIDRequired.Error(),
			},
		},
		{
			name:         "download file fails due to whitespace fileID",
			auth:         true,
			fileID:       "   ",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"fileID": modelUtils.ErrFileIDRequired.Error(),
			},
		},
		{
			name:         "download file fails due to invalid fileID",
			auth:         true,
			fileID:       "invalid-uuid",
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"fileID": modelUtils.ErrFileIDInvalidUUID.Error(),
			},
		},
		{
			name:          "download file fails due to missing metadata",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:   mocks.GrpcOpMissingMetadata,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "download file fails due to missing userID in metadata",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:   mocks.GrpcOpMissingUserID,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrUnauthorized.Error(),
		},
		{
			name:          "download file fails due to grpc internal error",
			auth:          true,
			fileID:        "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerErrors.ErrFailedToDownloadFile.Error(),
		},
		{
			name:         "download file succeeds when file not found",
			auth:         true,
			fileID:       "550e8400-e29b-41d4-a716-446655440000",
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedBody: config.NullString,
			expectedType: "application/octet-stream",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTest(http.MethodGet, "/api/files/:fileID/download", config.NullString, tt.auth)
			c.Params = []gin.Param{{Key: "fileID", Value: tt.fileID}}

			mmsClient := &mocks.MockMMSClient{DownloadGrpcErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{}
			client := grpc.NewClient(mmsClient, svc)
			handler := NewUserFilesHandler(client, svc)

			handler.DownloadFileHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedError != nil {
				mockUtils.CheckError(t, w, tt.expectedError)
			}

			if tt.expectedCode == http.StatusOK {
				if w.Body.String() != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, w.Body.String())
				}
				if ct := w.Header().Get("Content-Type"); ct != tt.expectedType {
					t.Errorf("expected content-type %q, got %q", tt.expectedType, ct)
				}
			}
		})
	}
}
