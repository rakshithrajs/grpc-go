package handlers

import (
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func TestRegisterUserHandler(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		mockErr       mocks.DbOperationError
		expectedCode  int
		expectedError any
		expectedData  any
	}{
		{
			name:         "user registration successful",
			body:         `{"name":"Test","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			mockErr:      mocks.DbOpSuccess,
			expectedCode: http.StatusCreated,
			expectedData: map[string]any{
				"user": map[string]any{
					"id":           "success-user-id",
					"name":         "Test",
					"email":        "test@example.com",
					"phone":        "1234567890",
					"createdAtUTC": "0001-01-01T00:00:00Z",
					"updatedAtUTC": "0001-01-01T00:00:00Z",
				},
			},
		},
		{
			name:          "user registration fails as of invalid json",
			body:          `{`,
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrInvalidJSON.Error(),
		},
		{
			name:         "user registration fails as of empty name",
			body:         `{"name":"","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"name": utils.ErrNameRequired.Error(),
			},
		},
		{
			name:         "user registration fails as of invalid name",
			body:         `{"name":"Test123","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"name": utils.ErrInvalidName.Error(),
			},
		},
		{
			name:         "user registration fails as of empty email",
			body:         `{"name":"Test","email":"","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"email": utils.ErrEmailRequired.Error(),
			},
		},
		{
			name:          "user registration fails as of invalid email format",
			body:          `{"name":"Test","email":"invalid-email","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			expectedCode:  http.StatusBadRequest,
			expectedError: map[string]string{"email": utils.ErrInvalidEmail.Error()},
		},
		{
			name:         "user registration fails as of empty password",
			body:         `{"name":"Test","email":"test@example.com","password":"","confirmPassword":"","phone":"1234567890"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password":        utils.ErrPasswordRequired.Error(),
				"confirmPassword": utils.ErrPasswordConfirmRequired.Error(),
			},
		},
		{
			name:         "user registration fails as of password mismatch",
			body:         `{"name":"Test","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"DifferentPassword@123","phone":"1234567890"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"confirmPassword": utils.ErrPasswordMismatch.Error(),
			},
		},
		{
			name:         "user registration fails as of empty phone",
			body:         `{"name":"Test","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":""}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"phone": utils.ErrPhoneRequired.Error(),
			},
		},
		{
			name:         "user registration fails as of invalid phone",
			body:         `{"name":"Test","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"123"}`,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"phone": utils.ErrInvalidPhoneNumber.Error(),
			},
		},
		{
			name:          "user registration fails as of duplicate email",
			body:          `{"name":"Test","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			mockErr:       mocks.DbOpDuplicateEmail,
			expectedCode:  http.StatusConflict,
			expectedError: utils.ErrUserEmailAlreadyExists.Error(),
		},
		{
			name:          "user registration fails as of duplicate phone",
			body:          `{"name":"Test","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			mockErr:       mocks.DbOpDuplicatePhone,
			expectedCode:  http.StatusConflict,
			expectedError: utils.ErrPhoneNumberAlreadyExists.Error(),
		},
		{
			name:          "user registration fails as of internal server error",
			body:          `{"name":"Test","email":"test@example.com","password":"ValidPassword@123","confirmPassword":"ValidPassword@123","phone":"1234567890"}`,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrFailedToCreateUser.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodPost, "/api/users/register", tt.body, false)

			svc := &mocks.MockUserService{MockErr: tt.mockErr}
			handler := NewUMSHandler(svc)

			handler.RegisterUserHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusCreated {
				mocks.CheckData(t, w, tt.expectedData)
				return
			}

			if tt.expectedError != nil {
				mocks.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
