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
	MockErr           DbOperationError
	updateRollbackErr DbOperationError
	DiskWriteFailure  bool
	ReturnOldName     bool
	UserID            string
	FileID            string
	Files             []*models.ListFileResponse
	updateCallCount   int
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

	switch m.MockErr {
	case DbOpDuplicateName:
		return nil, storage.ErrFileNameAlreadyExists
	case DbOpInternalError:
		return nil, storage.ErrFailedToUploadFile
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

func (m *MockFileService) GetFiles(ctx context.Context, userID string) ([]*models.ListFileResponse, error) {
	m.UserID = userID

	switch m.MockErr {
	case DbOpInternalError:
		return nil, storage.ErrFailedToGetFiles
	}

	if m.Files == nil {
		mt := models.MimeTypeTextPlain
		m.Files = []*models.ListFileResponse{
			{ID: "file-1", FileName: "file1.txt", FileSize: 100, MimeType: mt},
			{ID: "file-2", FileName: "file2.txt", FileSize: 200, MimeType: mt},
		}
	}
	return m.Files, nil
}

func (m *MockFileService) GetFileByID(ctx context.Context, id, userID string) (*models.File, error) {
	m.UserID = userID
	m.FileID = id

	switch m.MockErr {
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

	err := m.MockErr
	if m.DiskWriteFailure {
		err = m.updateRollbackErr
	}

	switch err {
	case DbOpInternalError:
		return nil, storage.ErrFailedToUpdateFile
	case DbOpNotFound:
		return nil, storage.ErrFileNotFound
	case DbOpDuplicateName:
		return nil, storage.ErrFileNameAlreadyExists
	}

	m.updateCallCount++

	name := "test.txt"
	path := filepath.Join(userStorageDir(userID), name)
	if req.Name != config.NullString {
		name = req.Name
		path = filepath.Join(userStorageDir(userID), name)
	}

	if m.ReturnOldName {
		if m.updateCallCount == 1 || m.updateCallCount >= 3 {
			oldName := "old.txt"
			name = oldName
			path = filepath.Join(userStorageDir(userID), oldName)
		}
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

	switch m.MockErr {
	case DbOpInternalError:
		return nil, storage.ErrFailedToDeleteFile
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
