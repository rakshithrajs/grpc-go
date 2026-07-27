package mocks

import (
	"context"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

type MockUserFilesService struct {
	CreateUserFileError DbOperationError
	DeleteUserFileError DbOperationError
	ListUserFilesError  DbOperationError
	UpdateUserFileError DbOperationError
	UpdateRollbackError DbOperationError
	ReturnEmptyList     bool
	UserID              string
	FileID              string
	FileName            string
	Files               []models.UserFiles
	updateCallCount     int
}

func (m *MockUserFilesService) CreateUserFile(ctx context.Context, userID, fileID, fileName string) error {
	m.UserID = userID
	m.FileID = fileID
	m.FileName = fileName

	switch m.CreateUserFileError {
	case DbOpDuplicateFile:
		return handlerErrors.ErrUserFileAlreadyExists
	case DbOpInternalError:
		return storage.ErrFailedToCreateUserFile
	case DbOpRollbackFailure:
		return handlerErrors.ErrFailedToRollback
	}

	return nil
}

func (m *MockUserFilesService) DeleteUserFile(ctx context.Context, userID, fileID string) (string, error) {
	m.UserID = userID
	m.FileID = fileID

	switch m.DeleteUserFileError {
	case DbOpInternalError:
		return config.NullString, storage.ErrFailedToDeleteUserFile
	case DbOpNotFound:
		return config.NullString, nil
	case DbOpRollbackFailure:
		return config.NullString, handlerErrors.ErrFailedToRollback
	}

	m.FileName = "test-file-name.txt"
	return m.FileName, nil
}

func (m *MockUserFilesService) ListUserFiles(ctx context.Context, userID string) ([]models.UserFiles, error) {
	m.UserID = userID

	switch m.ListUserFilesError {
	case DbOpInternalError:
		return nil, storage.ErrFailedToListUserFiles
	}

	if m.ReturnEmptyList {
		return []models.UserFiles{}, nil
	}

	m.Files = []models.UserFiles{
		{FileID: "file-1", UserID: userID, FileName: "file1.txt"},
		{FileID: "file-2", UserID: userID, FileName: "file2.txt"},
	}
	return m.Files, nil
}

func (m *MockUserFilesService) UpdateUserFile(ctx context.Context, userID, fileID, fileName string) (string, error) {
	m.UserID = userID
	m.FileID = fileID
	m.updateCallCount++

	if m.updateCallCount == 1 {
		switch m.UpdateUserFileError {
		case DbOpInternalError:
			return config.NullString, storage.ErrFailedToUpdateUserFile
		case DbOpNotFound:
			return config.NullString, nil
		}

		return "old-file-name.txt", nil
	}

	switch m.UpdateRollbackError {
	case DbOpInternalError:
		return config.NullString, handlerErrors.ErrFailedToRollback
	}

	return "old-file-name.txt", nil
}
