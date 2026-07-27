package mocks

import (
	"context"
	"path/filepath"
	"time"

	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"
)

type MockFileService struct {
	UploadFileErr     DbOperationError
	GetFileByIDErr    DbOperationError
	UpdateFileErr     DbOperationError
	UpdateRollbackErr DbOperationError
	DeleteFileErr DbOperationError
	ReturnOldName bool
	ReturnEmptyID bool
	UserID            string
	FileID            string
	Files             []*models.ListFileResponse
}

var ZeroTime = time.Time{}

func userStorageDir(userID string) string {
	cfg, err := config.GetConfig()
	if err != nil {
		return config.NullString
	}
	return filepath.Join(cfg.UserStoragePath, userID)
}

func (m *MockFileService) UploadFile(ctx context.Context, file *models.File) (*models.File, error) {
	m.UserID = file.UserID

	switch m.UploadFileErr {
	case DbOpDuplicateName:
		return nil, storage.ErrFileAlreadyExists
	case DbOpInternalError:
		return nil, storage.ErrFailedToUploadFile
	case DbOpRollbackError:
		return nil, storage.ErrFailedToRollback
	}

	return &models.File{
		ID:       "file-id-123",
		UserID:   file.UserID,
		Name:     file.Name,
		Path:     file.Path,
		Size:     file.Size,
		MimeType: file.MimeType,
	}, nil
}

func (m *MockFileService) GetFileByID(ctx context.Context, id, userID string) (*models.File, error) {
	m.UserID = userID
	m.FileID = id

	switch m.GetFileByIDErr {
	case DbOpInternalError:
		return nil, storage.ErrFailedToGetFileByID
	case DbOpNotFound:
		return nil, storage.ErrFileNotFound
	}

	name := "test.txt"
	return &models.File{
		ID:           id,
		UserID:       userID,
		Name:         name,
		Path:         filepath.Join(userStorageDir(userID), name),
		Size:         12,
		MimeType:     models.MimeTypeTextPlain,
		CreatedAtUTC: ZeroTime,
		UpdatedAtUTC: ZeroTime,
	}, nil
}

func (m *MockFileService) UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) (*models.File, error) {
	m.UserID = userID
	m.FileID = id

	activeErr := m.UpdateFileErr
	if req.Path != config.NullString && m.UpdateRollbackErr != DbOpSuccess {
		activeErr = m.UpdateRollbackErr
	}

	switch activeErr {
	case DbOpInternalError:
		return nil, storage.ErrFailedToUpdateFile
	case DbOpNotFound:
		return &models.File{}, nil
	case DbOpDuplicateName:
		return nil, storage.ErrFileAlreadyExists
	}

	if m.ReturnEmptyID {
		return &models.File{ID: config.NullString}, nil
	}

	originalName := "test.txt"
	if m.ReturnOldName {
		originalName = "old.txt"
	}

	name := originalName
	path := filepath.Join(userStorageDir(userID), name)

	if req.Name != config.NullString {
		name = req.Name
	}
	if req.Path != config.NullString {
		path = req.Path
	}

	return &models.File{
		ID:           id,
		UserID:       userID,
		Name:         name,
		Path:         path,
		Size:         12,
		MimeType:     models.MimeTypeTextPlain,
		CreatedAtUTC: ZeroTime,
		UpdatedAtUTC: ZeroTime,
	}, nil
}

func (m *MockFileService) DeleteFile(ctx context.Context, id, userID string) (*models.File, error) {
	m.UserID = userID
	m.FileID = id

	switch m.DeleteFileErr {
	case DbOpInternalError:
		return nil, storage.ErrFailedToDeleteFile
	case DbOpNotFound:
		return nil, nil
	}

	name := "test.txt"
	return &models.File{
		ID:           id,
		UserID:       userID,
		Name:         name,
		Path:         filepath.Join(userStorageDir(userID), name),
		Size:         12,
		MimeType:     models.MimeTypeTextPlain,
		CreatedAtUTC: ZeroTime,
		UpdatedAtUTC: ZeroTime,
	}, nil
}
