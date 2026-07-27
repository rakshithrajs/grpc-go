package storage

import (
	"context"

	"github.com/rakshithrajs/cloud/MMS/internal/models"
)

type FileService interface {
	UploadFile(ctx context.Context, file *models.File) (*models.File, error)
	GetFileByID(ctx context.Context, id string, userID string) (*models.File, error)
	UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) (*models.File, error)
	DeleteFile(ctx context.Context, id string, userID string) (*models.File, error)
}
