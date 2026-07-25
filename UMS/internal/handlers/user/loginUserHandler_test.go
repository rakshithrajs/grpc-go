package handlers

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
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
			name:          "user login fails as of invalid json",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrInvalidJSON.Error(),
		},
		{
			name:         "user login fails as of empty email",
			body:         `{"email":config.NullString,"password":"ValidPassword@123"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": modelUtils.ErrEmailRequired.Error(),
			},
		},
		{
			name:         "user login fails as of empty password",
			body:         `{"email":"test@example.com","password":config.NullString}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": modelUtils.ErrPasswordRequired.Error(),
			},
		},
		{
			name:         "user login fails as of empty email and password",
			body:         `{"email":config.NullString,"password":config.NullString}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email":    modelUtils.ErrEmailRequired.Error(),
				"password": modelUtils.ErrPasswordRequired.Error(),
			},
		},
		{
			name:          "user login fails as of invalid email format",
			body:          `{"email":"invalid-email","password":"ValidPassword@123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of invalid email domain",
			body:          `{"email":"test@invalid.com","password":"ValidPassword@123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of invalid password",
			body:          `{"email":"test@example.com","password":"123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of short password",
			body:          `{"email":"test@example.com","password":"12345"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "user login fails as of email does not exist",
			body:          `{"email":"test@example.com","password":"ValidPassword@123"}`,
			mockErr:       mocks.DbOpNotFound,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrInvalidCredentials.Error(),
		},
		{
			name:          "email exists but internal server error",
			body:          `{"email":"test@example.com","password":"ValidPassword@123"}`,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerUtils.ErrFailedToLoginUser.Error(),
		},
		{
			name:          "email exists but password does not match",
			body:          `{"email":"test@example.com","password":"WrongPassword@123"}`,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrInvalidCredentials.Error(),
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
