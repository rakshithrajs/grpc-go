package handlers

import (
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func TestGetUserProfileHandler(t *testing.T) {
	tests := []struct {
		name          string
		auth          bool
		mockErr       mocks.DbOperationError
		expectedCode  int
		expectedError any
		expectedData  any
	}{
		{
			name:          "get user profile fails due to missing auth",
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "get user profile fails due to internal server error",
			auth:          true,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrSomethingWentWrong.Error(),
		},
		{
			name:         "get user profile succeeds with no user found",
			auth:         true,
			mockErr:      mocks.DbOpNotFound,
			expectedCode: http.StatusOK,
			expectedData: map[string]any{
				"user": nil,
			},
		},
		{
			name:         "get user profile success",
			auth:         true,
			expectedCode: http.StatusOK,
			expectedData: map[string]any{
				"user": map[string]any{
					"id":           "test-user-id",
					"name":         "Test User",
					"email":        "test@example.com",
					"phone":        "1234567890",
					"createdAtUTC": "0001-01-01T00:00:00Z",
					"updatedAtUTC": "0001-01-01T00:00:00Z",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodGet, "/api/users/profile", "", tt.auth)

			svc := &mocks.MockUserService{GetUserByIDErr: tt.mockErr}
			handler := NewUserHandler(svc)

			handler.GetUserProfileHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				mocks.CheckData(t, w, tt.expectedData)
				return
			}

			if tt.expectedError != nil {
				mocks.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
