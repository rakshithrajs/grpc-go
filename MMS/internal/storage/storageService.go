package storage

import (
	"context"

	"github.com/rakshithrajs/cloud/MMS/internal/models"
)

type MMService interface {
	UploadFile(ctx context.Context, file *models.File) (*models.File, error)
	GetMMS(ctx context.Context, userID string) ([]*models.ListFileResponse, error)
	GetFileByID(ctx context.Context, id string, userID string) (*models.File, error)
	UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) error
	DeleteFile(ctx context.Context, id string, userID string) error
}
