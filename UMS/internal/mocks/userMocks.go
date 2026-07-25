package mocks

import (
	"context"
	"time"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

var mockPasswordHash string

var ZeroTime = time.Time{}

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

func (m *MockUserService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	switch m.CreateUserErr {
	case DbOpDuplicateEmail:
		return nil, utils.ErrUserEmailAlreadyExists
	case DbOpDuplicatePhone:
		return nil, utils.ErrPhoneNumberAlreadyExists
	case DbOpInternalError:
		return nil, utils.ErrFailedToCreateUser
	}

	m.User = user
	user.ID = "success-user-id"
	user.CreatedAtUTC = ZeroTime
	user.UpdatedAtUTC = ZeroTime
	user.Password = config.NullString
	return user, nil
}

func (m *MockUserService) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	m.ID = id

	switch m.GetUserByIDErr {
	case DbOpNotFound:
		return nil, utils.ErrUserNotFound
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

func (m *MockUserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	m.Email = email

	switch m.GetUserByEmailErr {
	case DbOpNotFound:
		return nil, utils.ErrEmailNotFound
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

func (m *MockUserService) UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error {
	m.ID = id
	m.UpdateReq = req

	switch m.UpdateUserErr {
	case DbOpDuplicateEmail:
		return utils.ErrUserEmailAlreadyExists
	case DbOpDuplicatePhone:
		return utils.ErrPhoneNumberAlreadyExists
	case DbOpInternalError:
		return utils.ErrFailedToUpdateUser
	}

	return nil
}

func init() {
	config.SetConfig(&config.Config{JWTSecret: "test-secret"})

	hash, err := bcrypt.GenerateFromPassword([]byte("ValidPassword@123"), bcrypt.MinCost)
	if err != nil {
		panic(err)
	}
	mockPasswordHash = string(hash)
}
