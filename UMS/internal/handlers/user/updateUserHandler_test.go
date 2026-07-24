package handlers

import (
	"net/http"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
)

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		auth          bool
		mockErr       mocks.DbOperationError
		expectedCode  int
		expectedError any
	}{
		{
			name:          "user update fails due to missing auth",
			body:          `{"phone":"0987654321"}`,
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: utils.ErrUnauthorized.Error(),
		},
		{
			name:          "user update fails due to invalid json",
			body:          `{`,
			auth:          true,
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrInvalidJSON.Error(),
		},
		{
			name:          "user update fails due to no fields to update",
			body:          `{}`,
			auth:          true,
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrNoFieldsToUpdate.Error(),
		},
		{
			name:         "user update fails due to invalid phone",
			body:         `{"phone":"123"}`,
			auth:         true,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"phone": utils.ErrInvalidPhoneNumber.Error(),
			},
		},
		{
			name:          "user update fails due to duplicate phone",
			body:          `{"phone":"0987654321"}`,
			auth:          true,
			mockErr:       mocks.DbOpDuplicatePhone,
			expectedCode:  http.StatusConflict,
			expectedError: utils.ErrPhoneNumberAlreadyExists.Error(),
		},
		{
			name:          "user update fails due to internal server error",
			body:          `{"phone":"0987654321"}`,
			auth:          true,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: utils.ErrFailedToUpdateUser.Error(),
		},
		{
			name:          "user update fails due to same old password",
			body:          `{"password":"ValidPassword@123","confirmPassword":"ValidPassword@123"}`,
			auth:          true,
			expectedCode:  http.StatusBadRequest,
			expectedError: utils.ErrPasswordSameAsOldPassword.Error(),
		},
		{
			name:         "user update fails due to invalid password",
			body:         `{"password":"123","confirmPassword":"123"}`,
			auth:         true,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": utils.ErrInvalidPassword.Error(),
			},
		},
		{
			name:         "user update success with phone",
			body:         `{"phone":"0987654321"}`,
			auth:         true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "user update success with password",
			body:         `{"password":"NewPassword@123","confirmPassword":"NewPassword@123"}`,
			auth:         true,
			expectedCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, w := mocks.SetUpGinTest(http.MethodPatch, "/api/users/update", tt.body, tt.auth)

			svc := &mocks.MockUserService{MockErr: tt.mockErr}
			handler := NewUserHandler(svc)

			handler.UpdateUserHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				mocks.CheckData(t, w, map[string]string{"message": userUpdatedMessage})
				return
			}

			if tt.expectedError != nil {
				mocks.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
