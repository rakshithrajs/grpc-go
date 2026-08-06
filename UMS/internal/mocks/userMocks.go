package mocks

import (
	"context"
	"time"

	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var mockPasswordHash string

// ZeroTime is a zero-value time used for deterministic mock timestamps.
var ZeroTime = time.Time{}

// MockUserService is a mock implementation of the user storage service.
type MockUserService struct {
	GetUserByIDErr    DbOperationError
	CreateUserErr     DbOperationError
	GetUserByEmailErr DbOperationError
	UpdateUserErr     DbOperationError
	User              *models.User
	ID                string
	Email             string
	UpdateReq         models.UpdateUserRequest
}

// CreateUser mocks inserting a new user into the database.
func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	switch m.CreateUserErr {
	case DbOpDuplicateEmail:
		return nil, handlerErrors.ErrUserEmailAlreadyExists
	case DbOpDuplicatePhone:
		return nil, handlerErrors.ErrPhoneNumberAlreadyExists
	case DbOpInternalError:
		return nil, handlerErrors.ErrFailedToCreateUser
	}

	m.User = &models.User{
		ID:           "success-user-id",
		Name:         user.Name,
		Email:        user.Email,
		Phone:        user.Phone,
		CreatedAtUTC: ZeroTime,
		UpdatedAtUTC: ZeroTime,
	}
	return m.User, nil
}

// GetUserByID mocks fetching a user by their ID from the database.
func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	m.ID = id

	switch m.GetUserByIDErr {
	case DbOpNotFound:
		return nil, nil
	case DbOpInternalError:
		return nil, storage.ErrFailedToGetUserByID
	}

	m.User = &models.User{
		ID:       id,
		Name:     "Test User",
		Email:    "test@example.com",
		Password: mockPasswordHash,
		Phone:    "1234567890",
	}
	return m.User, nil
}

// GetUserByEmail mocks fetching a user by their email from the database.
func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.Email = email

	switch m.GetUserByEmailErr {
	case DbOpNotFound:
		return nil, nil
	case DbOpInternalError:
		return nil, storage.ErrFailedToGetUserByEmail
	}

	m.User = &models.User{
		ID:       "test-user-id",
		Name:     "Test User",
		Email:    email,
		Password: mockPasswordHash,
		Phone:    "1234567890",
	}
	return m.User, nil
}

// UpdateUser mocks updating an existing user's details in the database.
func (m *MockUserService) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error {
	m.ID = id
	m.UpdateReq = req

	switch m.UpdateUserErr {
	case DbOpDuplicateEmail:
		return handlerErrors.ErrUserEmailAlreadyExists
	case DbOpDuplicatePhone:
		return handlerErrors.ErrPhoneNumberAlreadyExists
	case DbOpInternalError:
		return handlerErrors.ErrFailedToUpdateUser
	}

	return nil
}

func init() {
	hash, err := bcrypt.GenerateFromPassword([]byte("ValidPassword@123"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	mockPasswordHash = string(hash)
}
