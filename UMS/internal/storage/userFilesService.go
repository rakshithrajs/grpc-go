package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
)

const (
	fnCreateUserFile  = "CreateUserFile"
	fnDeleteUserFile  = "DeleteUserFile"
	fnListUserFiles   = "ListUserFiles"
	fnUpdateUserFile  = "UpdateUserFile"
	fnGetUserFileName = "GetUserFileName"
)

type userFilesStore struct {
	db *sql.DB
}

func NewUserFilesStore(db *sql.DB) UserFilesService {
	return &userFilesStore{db: db}
}

func (u *userFilesStore) CreateUserFile(ctx context.Context, userID, fileID, fileName string) error {
	query := `INSERT INTO "userFiles" ("userID", "fileID", "fileName") VALUES ($1, $2, $3)`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnCreateUserFile)+"prepare statement", slog.Any("error", err))
		return ErrFailedToCreateUserFile
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, userID, fileID, fileName); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == "userFiles_userID_fileID_unique" {
			return ErrUserFileAlreadyExists
		}
		slog.Error(logPrefix(fnCreateUserFile)+"execute statement", slog.Any("error", err))
		return ErrFailedToCreateUserFile
	}

	return nil
}

func (u *userFilesStore) DeleteUserFile(ctx context.Context, userID, fileID string) error {
	query := `DELETE FROM "userFiles" WHERE "userID" = $1 AND "fileID" = $2`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnDeleteUserFile)+"prepare statement", slog.Any("error", err))
		return ErrFailedToDeleteUserFile
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, userID, fileID)
	if err != nil {
		slog.Error(logPrefix(fnDeleteUserFile)+"execute statement", slog.Any("error", err))
		return ErrFailedToDeleteUserFile
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error(logPrefix(fnDeleteUserFile)+"getting rows affected", slog.Any("error", err))
		return ErrFailedToDeleteUserFile
	}

	if rowsAffected == 0 {
		slog.Warn(logPrefix(fnDeleteUserFile)+"no user file mapping found", slog.String("userID", userID), slog.String("fileID", fileID))
	}

	return nil
}

func (u *userFilesStore) ListUserFiles(ctx context.Context, userID string) ([]models.UserFiles, error) {
	query := `SELECT "fileID", "fileName" FROM "userFiles" WHERE "userID" = $1 ORDER BY "createdAtUTC" DESC`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"prepare statement", slog.Any("error", err))
		return nil, ErrFailedToListUserFiles
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"execute statement", slog.Any("error", err))
		return nil, ErrFailedToListUserFiles
	}
	defer rows.Close()

	files := []models.UserFiles{}
	for rows.Next() {
		var file models.UserFiles
		file.UserID = userID
		if err := rows.Scan(&file.FileID, &file.FileName); err != nil {
			slog.Error(logPrefix(fnListUserFiles)+"scan row", slog.Any("error", err))
			return nil, ErrFailedToListUserFiles
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"rows iteration", slog.Any("error", err))
		return nil, ErrFailedToListUserFiles
	}

	return files, nil
}

func (u *userFilesStore) UpdateUserFile(ctx context.Context, userID, fileID, fileName string) error {
	query := `UPDATE "userFiles" SET "fileName" = $1 WHERE "userID" = $2 AND "fileID" = $3`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnUpdateUserFile)+"prepare statement", slog.Any("error", err))
		return ErrFailedToUpdateUserFile
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, fileName, userID, fileID)
	if err != nil {
		slog.Error(logPrefix(fnUpdateUserFile)+"execute statement", slog.Any("error", err))
		return ErrFailedToUpdateUserFile
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error(logPrefix(fnUpdateUserFile)+"getting rows affected", slog.Any("error", err))
		return ErrFailedToUpdateUserFile
	}

	if rowsAffected == 0 {
		slog.Warn(logPrefix(fnUpdateUserFile)+"no user file mapping found", slog.String("userID", userID), slog.String("fileID", fileID))
	}

	return nil
}

func (u *userFilesStore) GetUserFileName(ctx context.Context, userID, fileID string) (string, error) {
	query := `SELECT "fileName" FROM "userFiles" WHERE "userID" = $1 AND "fileID" = $2`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnGetUserFileName)+"prepare statement", slog.Any("error", err))
		return "", ErrFailedToListUserFiles
	}
	defer stmt.Close()

	var fileName string
	if err := stmt.QueryRowContext(ctx, userID, fileID).Scan(&fileName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		slog.Error(logPrefix(fnGetUserFileName)+"execute statement", slog.Any("error", err))
		return "", ErrFailedToListUserFiles
	}

	return fileName, nil
}
