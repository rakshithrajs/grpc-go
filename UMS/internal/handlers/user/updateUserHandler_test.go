package handlers

import (
	"net/http"
	"testing"

	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	middlewareUtils "github.com/rakshithrajs/cloud/UMS/internal/middleware/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
)

func TestUpdateUserHandler(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		auth           bool
		mockErr        mocks.DbOperationError
		getUserByIDErr mocks.DbOperationError
		expectedCode   int
		expectedError  any
	}{
		{
			name:          "user update fails due to missing auth",
			body:          `{"phone":"0987654321"}`,
			auth:          false,
			expectedCode:  http.StatusUnauthorized,
			expectedError: handlerUtils.ErrUnauthorized.Error(),
		},
		{
			name:          "user update fails due to invalid json",
			body:          `{`,
			auth:          true,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrInvalidJSON.Error(),
		},
		{
			name:          "user update fails due to no fields to update",
			body:          `{}`,
			auth:          true,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrNoFieldsToUpdate.Error(),
		},
		{
			name:         "user update fails due to invalid phone",
			body:         `{"phone":"123"}`,
			auth:         true,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"phone": modelUtils.ErrInvalidPhoneNumber.Error(),
			},
		},
		{
			name:          "user update fails due to duplicate phone",
			body:          `{"phone":"0987654321"}`,
			auth:          true,
			mockErr:       mocks.DbOpDuplicatePhone,
			expectedCode:  http.StatusConflict,
			expectedError: handlerUtils.ErrPhoneNumberAlreadyExists.Error(),
		},
		{
			name:          "user update fails due to duplicate email",
			body:          `{"email":"new@example.com"}`,
			auth:          true,
			mockErr:       mocks.DbOpDuplicateEmail,
			expectedCode:  http.StatusConflict,
			expectedError: handlerUtils.ErrUserEmailAlreadyExists.Error(),
		},
		{
			name:           "user update fails due get user by ID error",
			body:           `{"password":"NewPassword@123","confirmPassword":"NewPassword@123"}`,
			auth:           true,
			mockErr:        mocks.DbOpSuccess,
			getUserByIDErr: mocks.DbOpInternalError,
			expectedCode:   http.StatusInternalServerError,
			expectedError:  middlewareUtils.ErrSomethingWentWrong.Error(),
		},
		{
			name:          "user update fails due to internal server error",
			body:          `{"phone":"0987654321"}`,
			auth:          true,
			mockErr:       mocks.DbOpInternalError,
			expectedCode:  http.StatusInternalServerError,
			expectedError: handlerUtils.ErrFailedToUpdateUser.Error(),
		},
		{
			name:          "user update fails due to same old password",
			body:          `{"password":"ValidPassword@123","confirmPassword":"ValidPassword@123"}`,
			auth:          true,
			expectedCode:  http.StatusBadRequest,
			expectedError: handlerUtils.ErrPasswordSameAsOldPassword.Error(),
		},
		{
			name:         "user update fails due to invalid password",
			body:         `{"password":"123","confirmPassword":"123"}`,
			auth:         true,
			expectedCode: http.StatusBadRequest,
			expectedError: map[string]string{
				"password": modelUtils.ErrInvalidPassword.Error(),
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
			c, w := mockUtils.SetUpGinTest(http.MethodPatch, "/api/users/update", tt.body, tt.auth)

			svc := &mocks.MockUserService{UpdateUserErr: tt.mockErr, GetUserByIDErr: tt.getUserByIDErr}
			handler := NewUserHandler(svc)

			handler.UpdateUserHandler(c)

			if w.Code != tt.expectedCode {
				t.Errorf("expected %d, got %d", tt.expectedCode, w.Code)
			}

			if tt.expectedCode == http.StatusOK {
				mockUtils.CheckData(t, w, map[string]string{"message": userUpdatedMessage})
				return
			}

			if tt.expectedError != nil {
				mockUtils.CheckError(t, w, tt.expectedError)
			}
		})
	}
}
