package storage

import (
	"context"

	"github.com/rakshithrajs/cloud/MMS/internal/models"
)

// FileService defines the contract for file persistence operations.
type FileService interface {
	CreateFile(ctx context.Context, file *models.File) (*models.File, error)
	GetFileByID(ctx context.Context, fileID string, userID string) (*models.File, error)
	UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) (*models.File, error)
	DeleteFile(ctx context.Context, id string, userID string) (*models.File, error)
}
