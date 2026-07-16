package storage

import (
	"cloud/internal/models"
	"context"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, id string, req models.UpdateUserRequest) error
}

type FileService interface {
	UploadFile(ctx context.Context, file *models.File) (*models.File, error)
	GetFiles(ctx context.Context, userID string) ([]*models.ListFileResponse, error)
	GetFileByID(ctx context.Context, id string, userID string) (*models.File, error)
	UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) error
	DeleteFile(ctx context.Context, id string, userID string) error
}
