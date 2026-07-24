package mocks

import (
	"context"

	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
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
	switch m.CreateUserFileError {
	case DbOpDuplicateFile:
		return utils.ErrUserFileAlreadyExists
	case DbOpInternalError:
		return storage.ErrFailedToCreateUserFile
	}

	m.UserID = userID
	m.FileID = fileID
	m.FileName = fileName
	return nil
}

func (m *MockUserFilesService) DeleteUserFile(ctx context.Context, userID, fileID string) (string, error) {
	m.UserID = userID
	m.FileID = fileID

	switch m.DeleteUserFileError {
	case DbOpInternalError:
		return "", storage.ErrFailedToDeleteUserFile
	case DbOpNotFound:
		return "", nil
	}

	return "test-file-name.txt", nil
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
			return "", storage.ErrFailedToUpdateUserFile
		case DbOpNotFound:
			return "", nil
		}

		oldName := "old-file-name.txt"
		m.FileName = fileName
		return oldName, nil
	}

	switch m.UpdateRollbackError {
	case DbOpInternalError:
		return "", utils.ErrFailedToRollback
	}

	return "old-file-name.txt", nil
}
