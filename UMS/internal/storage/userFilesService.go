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
	// function name for CreateUserFile
	fnCreateUserFile = "CreateUserFile"

	// function name for DeleteUserFile
	fnDeleteUserFile = "DeleteUserFile"

	// function name for ListUserFiles
	fnListUserFiles = "ListUserFiles"

	// function name for UpdateUserFile
	fnUpdateUserFile = "UpdateUserFile"

	// primary key constraint name for userFiles table on userID and fileID columns
	primaryKeyConstraintUserFiles = "userFiles_pkey"
)

type userFilesStore struct {
	db *sql.DB
}

// NewUserFilesStore creates a new PostgreSQL-backed UserFilesService.
func NewUserFilesStore(db *sql.DB) UserFilesService {
	return &userFilesStore{db: db}
}

// CreateUserFile persists a mapping between a user and an uploaded file.
func (u *userFilesStore) CreateUserFile(ctx context.Context, userID, fileID, fileName string) error {
	query := `INSERT INTO "userFiles" ("userID", "fileID", "fileName") VALUES ($1, $2, $3)`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnCreateUserFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return ErrFailedToCreateUserFile
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnCreateUserFile)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

	if _, err := stmt.ExecContext(ctx, userID, fileID, fileName); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == primaryKeyConstraintUserFiles {
			return handlerErrors.ErrUserFileAlreadyExists
		}
		slog.Error(logPrefix(fnCreateUserFile)+"execute statement", slog.Any(config.ErrorKey, err))
		return ErrFailedToCreateUserFile
	}

	return nil
}

// DeleteUserFile removes a user-file mapping and returns the deleted filename.
func (u *userFilesStore) DeleteUserFile(ctx context.Context, userID, fileID string) (string, error) {
	query := `DELETE FROM "userFiles" WHERE "userID" = $1 AND "fileID" = $2 RETURNING "fileName"`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnDeleteUserFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return config.NullString, ErrFailedToDeleteUserFile
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnDeleteUserFile)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

	var fileName string
	if err := stmt.QueryRowContext(ctx, userID, fileID).Scan(&fileName); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fileName, nil
		}
		slog.Error(logPrefix(fnDeleteUserFile)+"query row", slog.Any(config.ErrorKey, err))
		return config.NullString, ErrFailedToDeleteUserFile
	}

	return fileName, nil
}

// ListUserFiles returns all files owned by the given user.
func (u *userFilesStore) ListUserFiles(ctx context.Context, userID string) ([]models.UserFiles, error) {
	query := `SELECT "fileID", "fileName" FROM "userFiles" WHERE "userID" = $1 ORDER BY "createdAtUTC" DESC`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToListUserFiles
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnListUserFiles)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		slog.Error(logPrefix(fnListUserFiles)+"execute statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToListUserFiles
	}
	defer func() {
		if err := rows.Close(); err != nil {
			slog.Error(logPrefix(fnListUserFiles)+"failed to close rows", slog.Any(config.ErrorKey, err))
		}
	}()

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

// UpdateUserFile renames a user-file mapping and returns the previous filename.
func (u *userFilesStore) UpdateUserFile(ctx context.Context, userID, fileID, fileName string) (string, error) {
	query := `
		WITH old AS (
			SELECT "fileName" FROM "userFiles" WHERE "userID" = $2 AND "fileID" = $3
		)
		UPDATE "userFiles"
		SET "fileName" = $1, "updatedAtUTC" = NOW()
		FROM old
		WHERE "userFiles"."userID" = $2 AND "userFiles"."fileID" = $3
		RETURNING old."fileName"
	`

	stmt, err := u.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnUpdateUserFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return config.NullString, ErrFailedToUpdateUserFile
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnUpdateUserFile)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

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
