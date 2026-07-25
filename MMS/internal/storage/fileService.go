package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
	storageUtils "github.com/rakshithrajs/cloud/MMS/internal/storage/utils"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

const (
	fnUploadFile  = "UploadFile"
	fnGetFiles    = "GetFiles"
	fnGetFileByID = "GetFileByID"
	fnUpdateFile  = "UpdateFile"
	fnDeleteFile  = "DeleteFile"

	uniqueConstraintFilesUserName = "files_user_name_unique"
)

func logPrefix(fn string) string { return "[" + fn + "]: " }

type FileStore struct {
	db *sql.DB
}

func NewFileStore(db *sql.DB) FileService {
	return &FileStore{db: db}
}

func (f *FileStore) UploadFile(ctx context.Context, file *models.File) (*models.File, error) {
	query := `INSERT INTO files ("userID", name, path, size, "mimeType") VALUES ($1, $2, $3, $4, $5) RETURNING "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC"`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnUploadFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUploadFile
	}
	defer stmt.Close()

	var newFile models.File
	if err := stmt.QueryRowContext(ctx, file.UserID, file.Name, file.Path, file.Size, file.MimeType).Scan(&newFile.ID, &newFile.UserID, &newFile.Name, &newFile.Path, &newFile.Size, &newFile.MimeType, &newFile.CreatedAtUTC, &newFile.UpdatedAtUTC); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == uniqueConstraintFilesUserName {
			return nil, ErrFileNameAlreadyExists
		}
		slog.Error(logPrefix(fnUploadFile)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUploadFile
	}

	return &newFile, nil
}

func (f *FileStore) GetFiles(ctx context.Context, userID string) ([]*models.ListFileResponse, error) {
	query := `SELECT "ID", name, size, "mimeType" FROM files WHERE "userID" = $1`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnGetFiles)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetFiles
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		slog.Error(logPrefix(fnGetFiles)+"query rows", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetFiles
	}
	defer rows.Close()

	var files []*models.ListFileResponse
	for rows.Next() {
		var file models.ListFileResponse
		if err := rows.Scan(&file.ID, &file.FileName, &file.FileSize, &file.MimeType); err != nil {
			slog.Error(logPrefix(fnGetFiles)+"scan row", slog.Any(config.ErrorKey, err))
			return nil, ErrFailedToGetFiles
		}
		files = append(files, &file)
	}

	if err := rows.Err(); err != nil {
		slog.Error(logPrefix(fnGetFiles)+"rows error", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetFiles
	}

	return files, nil
}

func (f *FileStore) GetFileByID(ctx context.Context, id string, userID string) (*models.File, error) {
	query := `SELECT "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC" FROM files WHERE "ID" = $1 AND "userID" = $2`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnGetFileByID)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetFileByID
	}
	defer stmt.Close()

	var file models.File
	if err := stmt.QueryRowContext(ctx, id, userID).Scan(
		&file.ID, &file.UserID, &file.Name, &file.Path, &file.Size, &file.MimeType, &file.CreatedAtUTC, &file.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFileNotFound
		}
		slog.Error(logPrefix(fnGetFileByID)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetFileByID
	}

	return &file, nil
}

func (f *FileStore) UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) (*models.File, error) {
	fields := make([]storageUtils.UpdateField, 0, 2)
	if req.Name != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "name", Value: req.Name})
	}
	if req.Path != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "path", Value: req.Path})
	}

	query, args := storageUtils.BuildUpdateSQL("files", fields, []string{"ID", "userID"})
	query += ` RETURNING "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC"`
	args[0] = id
	args[1] = userID

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnUpdateFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUpdateFile
	}
	defer stmt.Close()

	var file models.File
	if err := stmt.QueryRowContext(ctx, args...).Scan(
		&file.ID, &file.UserID, &file.Name, &file.Path, &file.Size, &file.MimeType, &file.CreatedAtUTC, &file.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFileNotFound
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == uniqueConstraintFilesUserName {
			return nil, ErrFileNameAlreadyExists
		}
		slog.Error(logPrefix(fnUpdateFile)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUpdateFile
	}

	return &file, nil
}

func (f *FileStore) DeleteFile(ctx context.Context, id string, userID string) (*models.File, error) {
	query := `DELETE FROM files WHERE "ID" = $1 AND "userID" = $2 RETURNING "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC"`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnDeleteFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToDeleteFile
	}
	defer stmt.Close()

	var file models.File
	if err := stmt.QueryRowContext(ctx, id, userID).Scan(
		&file.ID, &file.UserID, &file.Name, &file.Path, &file.Size, &file.MimeType, &file.CreatedAtUTC, &file.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFileNotFound
		}
		slog.Error(logPrefix(fnDeleteFile)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToDeleteFile
	}

	return &file, nil
}
