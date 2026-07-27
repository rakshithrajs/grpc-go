# a usual db service in the project

```go
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
	// function name for UploadFile
	fnUploadFile = "UploadFile"

	// function name for GetFileByID
	fnGetFileByID = "GetFileByID"

	// function name for UpdateFile
	fnUpdateFile = "UpdateFile"

	// function name for DeleteFile
	fnDeleteFile = "DeleteFile"

	// unique constraint name for files table on userID and name columns
	uniqueConstraintFilesUserName = "files_user_name_unique"
)

// logPrefix returns a formatted string for logging purposes, including the function name.
func logPrefix(fn string) string { return "[" + fn + "]: " }

type FileStore struct {
	db *sql.DB
}

// NewFileStore creates a new instance of FileStore with the provided database connection.
func NewFileStore(db *sql.DB) FileService {
	return &FileStore{db: db}
}

func (f *FileStore) UploadFile(ctx context.Context, file *models.File) (*models.File, error) {
	query := `INSERT INTO files ("userID", name, path, size, "mimeType") VALUES ($1, $2, $3, $4, $5) RETURNING "ID", name, size, "mimeType"`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnUploadFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUploadFile
	}
	defer stmt.Close()

	var newFile models.File
	if err := stmt.QueryRowContext(ctx, file.UserID, file.Name, file.Path, file.Size, file.MimeType).Scan(&newFile.ID, &newFile.Name, &newFile.Size, &newFile.MimeType); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == uniqueConstraintFilesUserName {
			return nil, ErrFileAlreadyExists
		}
		slog.Error(logPrefix(fnUploadFile)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUploadFile
	}

	return &newFile, nil
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
			return &file, nil
		}
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == uniqueConstraintFilesUserName {
			return nil, ErrFileAlreadyExists
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
			return nil, nil
		}
		slog.Error(logPrefix(fnDeleteFile)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToDeleteFile
	}

	return &file, nil
}
```
