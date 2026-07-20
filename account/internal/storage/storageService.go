package storage

import (
	"github.com/rakshithrajs/cloud/services/account/internal/models"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error
}
