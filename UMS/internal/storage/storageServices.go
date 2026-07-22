package storage

import (
	"context"

	"github.com/rakshithrajs/cloud/UMS/internal/models"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error
}

type UserFilesService interface {
	CreateUserFile(ctx context.Context, userID, fileID, fileName string) error
	DeleteUserFile(ctx context.Context, userID, fileID string) error
	ListUserFiles(ctx context.Context, userID string) ([]models.UserFiles, error)
	UpdateUserFile(ctx context.Context, userID, fileID, fileName string) error
	GetUserFileName(ctx context.Context, userID, fileID string) (string, error)
}
