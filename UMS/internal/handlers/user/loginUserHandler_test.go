package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func TestLoginUserHandler(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		mockErr       mocks.DbOperationError
		expectedCode  int
		expectedError any
	}{
		{
			name:         "user login successful",
			body:         `{"email":"test@example.com","password":"ValidPassword@123"}`,
			mockErr:      mocks.DbOpSuccess,
			expectedCode: http.StatusOK,
		},
		{
			name:          "user login fails as of invalid json",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrInvalidJSON.Error(),
		},
		{
			name:         "user login fails as of empty email",
			body:         `{"email":"","password":"ValidPassword@123"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": utils.ErrEmailRequired.Error(),
			},
		},
		{
			name:         "user login fails as of empty password",
			body:         `{"email":"test@example.com","password":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": utils.ErrPasswordRequired.Error(),
			},
		},
		{
			name:         "user login fails as of empty email and password",
			body:         `{"email":"","password":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email":    utils.ErrEmailRequired.Error(),
				"password": utils.ErrPasswordRequired.Error(),
			},
		},
		{
			name:          "user login fails as of invalid email format",
			body:          `{"email":"invalid-email","password":"ValidPassword@123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of invalid email domain",
			body:          `{"email":"test@invalid.com","password":"ValidPassword@123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of invalid password",
			body:          `{"email":"test@example.com","password":"123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of short password",
			body:          `{"email":"test@example.com","password":"12345"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of email does not exist",
			body:          `{"email":"test@example.com","password":"ValidPassword@123"}`,
			mockErr:       mocks.DbOpNotFound,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "email exists but internal server error",
			body:          `{"email":"test@example.com","password":"ValidPassword@123"}`,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrFailedToLoginUser.Error(),
		},
		{
			name:          "email exists but password does not match",
			body:          `{"email":"test@example.com","password":"WrongPassword@123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrInvalidCredentials.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodPost, "/api/users/login", tt.body, false)

			svc := &mocks.MockUserService{MockErr: tt.mockErr}
			handler := NewUMSHandler(svc)

			handler.LoginUserHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				resp := make(map[string]string)
				json.NewDecoder(w.Body).Decode(&resp)
				if resp["token"] == "" {
					t.Errorf("expected token in response, got empty")
				}
			}

			if tt.expectedError != nil {
				mocks.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
