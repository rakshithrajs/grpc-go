package handlers

import (
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/grpc"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func TestListFilesHandler(t *testing.T) {
	tests := []struct {
		name          string
		auth          bool
		mockDbErr     mocks.DbOperationError
		returnEmpty   bool
		expectedCode  int
		expectedError any
		expectedData  any
	}{
		{
			name:          "list files fails due to missing auth",
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "list files fails due to db internal error",
			auth:          true,
			mockDbErr:     mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrFailedToListFiles.Error(),
		},
		{
			name:         "list files returns empty list",
			auth:         true,
			returnEmpty:  true,
			expectedCode: http.StatusOK,
			expectedData: map[string]any{"files": []any{}},
		},
		{
			name:         "list files succeeds",
			auth:         true,
			expectedCode: http.StatusOK,
			expectedData: map[string]any{
				"files": []any{
					map[string]any{"fileID": "file-1", "fileName": "file1.txt", "userID": "test-user-id"},
					map[string]any{"fileID": "file-2", "fileName": "file2.txt", "userID": "test-user-id"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodGet, "/api/files", "", tt.auth)

			mmsClient := &mocks.MockMMSClient{}
			svc := &mocks.MockUserFilesService{
				ListUserFilesError: tt.mockDbErr,
				ReturnEmptyList:    tt.returnEmpty,
			}
			client := grpc.NewClient(mmsClient, svc)
			handler := NewUserFilesHandler(client, svc)

			handler.ListFilesHandler(c)

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
