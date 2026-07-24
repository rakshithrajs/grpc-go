package handlers

import (
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
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
			name:          "download file fails due to missing auth",
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "download file fails due to missing fileID",
			auth:          true,
			fileID:        "",
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrFileIDRequired.Error(),
		},
		{
			name:          "download file fails due to missing metadata",
			auth:          true,
			fileID:        "file-id-123",
			mockGrpcErr:   mocks.GrpcOpMissingMetadata,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "download file fails due to missing userID in metadata",
			auth:          true,
			fileID:        "file-id-123",
			mockGrpcErr:   mocks.GrpcOpMissingUserID,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "download file fails due to grpc internal error",
			auth:          true,
			fileID:        "file-id-123",
			mockGrpcErr:   mocks.GrpcOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrFailedToDownloadFile.Error(),
		},
		{
			name:         "download file succeeds when file not found",
			auth:         true,
			fileID:       "file-id-123",
			mockGrpcErr:  mocks.GrpcOpNotFound,
			expectedCode: http.StatusOK,
			expectedBody: "",
			expectedType: "application/octet-stream",
		},
		{
			name:         "download file succeeds",
			auth:         true,
			fileID:       "file-id-123",
			expectedCode: http.StatusOK,
			expectedBody: "test content",
			expectedType: "text/plain",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodGet, "/api/files/:fileID/download", "", tt.auth)
			c.Params = []gin.Param{{Key: "fileID", Value: tt.fileID}}

			mmsClient := &mocks.MockMMSClient{MockErr: tt.mockGrpcErr}
			svc := &mocks.MockUserFilesService{}
			client := grpc.NewClient(mmsClient, svc)
			handler := NewUserFilesHandler(client, svc)

			handler.DownloadFileHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedError != nil {
				mocks.CheckError(t, w, tt.expectedError)
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
