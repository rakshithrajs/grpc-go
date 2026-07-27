package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
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
			name:          "user login fails due to invalid json",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerErrors.ErrInvalidJSON.Error(),
		},
		{
			name:         "user login fails due to empty email",
			body:         `{"email":"","password":"ValidPassword@123"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": modelUtils.ErrEmailRequired.Error(),
			},
		},
		{
			name:         "user login fails due to empty password",
			body:         `{"email":"test@example.com","password":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": modelUtils.ErrPasswordRequired.Error(),
			},
		},
		{
			name:         "user login fails due to empty email and password",
			body:         `{"email":"","password":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email":    modelUtils.ErrEmailRequired.Error(),
				"password": modelUtils.ErrPasswordRequired.Error(),
			},
		},
		{
			name:          "user login fails due to invalid email format",
			body:          `{"email":"invalid-email","password":"ValidPassword@123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: map[string]string{"email": modelUtils.ErrInvalidEmail.Error()},
		},
		{
			name:          "user login fails due to invalid email domain",
			body:          `{"email":"test@invalid.com","password":"ValidPassword@123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: map[string]string{"email": modelUtils.ErrInvalidEmail.Error()},
		},
		{
			name:          "user login fails due to invalid password",
			body:          `{"email":"test@example.com","password":"123"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: map[string]string{"password": modelUtils.ErrInvalidPassword.Error()},
		},
		{
			name:          "user login fails due to short password",
			body:          `{"email":"test@example.com","password":"12345"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: map[string]string{"password": modelUtils.ErrInvalidPassword.Error()},
		},
		{
			name:          "user login fails because email does not exist",
			body:          `{"email":"test@example.com","password":"ValidPassword@123"}`,
			mockErr:       mocks.DbOpNotFound,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrInvalidCredentials.Error(),
		},
		{
			name:          "email exists but internal server error",
			body:          `{"email":"test@example.com","password":"ValidPassword@123"}`,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerErrors.ErrFailedToLoginUser.Error(),
		},
		{
			name:          "email exists but password does not match",
			body:          `{"email":"test@example.com","password":"WrongPassword@123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerErrors.ErrInvalidCredentials.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mockUtils.SetUpGinTest(http.MethodPost, "/api/users/login", tt.body, false)

			svc := &mocks.MockUserService{GetUserByEmailErr: tt.mockErr}
			handler := NewUserHandler(svc)

			handler.LoginUserHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				resp := make(map[string]string)
				json.NewDecoder(w.Body).Decode(&resp)
				if resp["token"] == config.NullString {
					t.Errorf("expected token in response, got empty")
				}
			}

			if tt.expectedError != nil {
				mockUtils.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
