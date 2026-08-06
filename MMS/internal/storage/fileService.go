package storage

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"strings"

	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
	storageUtils "github.com/rakshithrajs/cloud/MMS/internal/storage/utils"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

const (
	// function name for CreateFile
	fnCreateFile = "CreateFile"

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

// FileStore implements FileService using a PostgreSQL database.
type FileStore struct {
	db *sql.DB
}

// NewFileStore creates a new instance of FileStore with the provided database connection.
func NewFileStore(db *sql.DB) FileService {
	return &FileStore{db: db}
}

// CreateFile inserts a new file record into the database.
func (f *FileStore) CreateFile(ctx context.Context, file *models.File) (*models.File, error) {
	query := `INSERT INTO files ("userID", "name", "path", "size", "mimeType") VALUES ($1, $2, $3, $4, $5) RETURNING "ID", "name", "size", "mimeType"`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnCreateFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUploadFile
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnCreateFile)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

	var newFile models.File
	if err := stmt.QueryRowContext(ctx, file.UserID, file.Name, file.Path, file.Size, file.MimeType).Scan(&newFile.ID, &newFile.Name, &newFile.Size, &newFile.MimeType); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation && pqErr.Constraint == uniqueConstraintFilesUserName {
			return nil, ErrFileAlreadyExists
		}
		slog.Error(logPrefix(fnCreateFile)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUploadFile
	}

	return &newFile, nil
}

// GetFileByID returns a file record by its ID and owner.
func (f *FileStore) GetFileByID(ctx context.Context, fileID string, userID string) (*models.File, error) {
	query := `SELECT "ID", "name", "path", "size", "mimeType" FROM files WHERE "ID" = $1 AND "userID" = $2`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnGetFileByID)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetFileByID
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnGetFileByID)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

	var file models.File
	file.UserID = userID
	if err := stmt.QueryRowContext(ctx, fileID, userID).Scan(
		&file.ID, &file.Name, &file.Path, &file.Size, &file.MimeType); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		slog.Error(logPrefix(fnGetFileByID)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToGetFileByID
	}

	return &file, nil
}

// UpdateFile applies a partial rename/move update to a file record and returns the old name and path.
func (f *FileStore) UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) (*models.File, error) {
	fields := make([]storageUtils.UpdateField, 0, 2)
	if req.Name != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "name", Value: req.Name})
	}
	if req.Path != config.NullString {
		fields = append(fields, storageUtils.UpdateField{Column: "path", Value: req.Path})
	}

	query, args := storageUtils.BuildUpdateSQL("files", fields, []string{"ID", "userID"})

	var queryBuilder strings.Builder
	queryBuilder.WriteString(`WITH old AS (SELECT "name", "path" FROM files WHERE "ID" = $1 AND "userID" = $2) `)
	queryBuilder.WriteString(query)
	queryBuilder.WriteString(` RETURNING (SELECT "name" FROM old), (SELECT "path" FROM old)`)

	args[0] = id
	args[1] = userID

	stmt, err := f.db.PrepareContext(ctx, queryBuilder.String())
	if err != nil {
		slog.Error(logPrefix(fnUpdateFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToUpdateFile
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnUpdateFile)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

	var file models.File
	file.UserID = userID
	if err := stmt.QueryRowContext(ctx, args...).Scan(
		&file.Name, &file.Path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
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

// DeleteFile removes a file record from the database.
func (f *FileStore) DeleteFile(ctx context.Context, id string, userID string) (*models.File, error) {
	query := `DELETE FROM files WHERE "ID" = $1 AND "userID" = $2 RETURNING "name", "path"`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error(logPrefix(fnDeleteFile)+"prepare statement", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToDeleteFile
	}
	defer func() {
		if err := stmt.Close(); err != nil {
			slog.Error(logPrefix(fnDeleteFile)+"failed to close statement", slog.Any(config.ErrorKey, err))
		}
	}()

	var file models.File
	file.UserID = userID
	if err := stmt.QueryRowContext(ctx, id, userID).Scan(
		&file.Name, &file.Path); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		slog.Error(logPrefix(fnDeleteFile)+"query row", slog.Any(config.ErrorKey, err))
		return nil, ErrFailedToDeleteFile
	}

	return &file, nil
}
