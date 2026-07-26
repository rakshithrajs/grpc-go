package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
)

const (
	fnCreateUserFile = "CreateUserFile"
	fnDeleteUserFile = "DeleteUserFile"
	fnListUserFiles  = "ListUserFiles"
	fnUpdateUserFile = "UpdateUserFile"

	uniqueConstraintUserFilesUserIDFileID = "userFiles_pkey"
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
		slog.Error(logPrefix(fnCreateUserFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return ErrFailedToCreateUserFile
	}
	defer stmt.Close()

	if _, err := stmt.ExecContext(ctx, userID, fileID, fileName); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == uniqueConstraintUserFilesUserIDFileID {
			return handlerErrors.ErrUserFileAlreadyExists
		}
		slog.Error(logPrefix(fnCreateUserFile)+"execute statement", slog.Any(config.ErrorKey, err))
		return ErrFailedToCreateUserFile
	}

	return nil
}

func (u *userFilesStore) DeleteUserFile(ctx context.Context, userID, fileID string) (string, error) {
	query := `DELETE FROM "userFiles" WHERE "userID" = $1 AND "fileID" = $2 RETURNING "fileName"`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnDeleteUserFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return config.NullString, ErrFailedToDeleteUserFile
	}
	defer stmt.Close()

	var fileName string
	if err := stmt.QueryRowContext(ctx, userID, fileID).Scan(&fileName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return config.NullString, nil
		}
		slog.Error(logPrefix(fnDeleteUserFile)+"query row", slog.Any(config.ErrorKey, err))
		return config.NullString, ErrFailedToDeleteUserFile
	}

	return fileName, nil
}

func (u *userFilesStore) ListUserFiles(ctx context.Context, userID string) ([]models.UserFiles, error) {
	query := `SELECT "fileID", "fileName" FROM "userFiles" WHERE "userID" = $1 ORDER BY "createdAtUTC" DESC`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToListUserFiles
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"execute statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToListUserFiles
	}
	defer rows.Close()

	files := []models.UserFiles{}
	for rows.Next() {
		var file models.UserFiles
		file.UserID = userID
		if err := rows.Scan(&file.FileID, &file.FileName); err != nil {
			slog.Error(logPrefix(fnListUserFiles)+"scan row", slog.Any(config.ErrorKey, err))
			return nil, ErrFailedToListUserFiles
		}
		files = append(files, file)
	}

	if err := rows.Err(); err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"rows iteration", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToListUserFiles
	}

	return files, nil
}

func (u *userFilesStore) UpdateUserFile(ctx context.Context, userID, fileID, fileName string) (string, error) {
	query := `UPDATE "userFiles" SET "fileName" = $1, "updatedAtUTC" = NOW() WHERE "userID" = $2 AND "fileID" = $3 RETURNING "fileName"`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnUpdateUserFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return config.NullString, ErrFailedToUpdateUserFile
	}
	defer stmt.Close()

	var oldFileName string
	if err := stmt.QueryRowContext(ctx, fileName, userID, fileID).Scan(&oldFileName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return config.NullString, nil
		}
		slog.Error(logPrefix(fnUpdateUserFile)+"query row", slog.Any(config.ErrorKey, err))
		return config.NullString, ErrFailedToUpdateUserFile
	}

	return oldFileName, nil
}
