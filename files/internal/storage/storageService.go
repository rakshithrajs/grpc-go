package storage

import (
	"context"

	"github.com/rakshithrajs/cloud/services/files/internal/models"
)

type FileService interface {
	UploadFile(ctx context.Context, file *models.File) (*models.File, error)
	GetFiles(ctx context.Context, userID string) ([]*models.ListFileResponse, error)
	GetFileByID(ctx context.Context, id string, userID string) (*models.File, error)
	UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) error
	DeleteFile(ctx context.Context, id string, userID string) error
}
