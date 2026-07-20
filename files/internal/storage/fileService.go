package storage

import (
	"github.com/rakshithrajs/cloud/services/files/internal/models"
	"github.com/rakshithrajs/cloud/services/files/internal/utils"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

type fileStore struct {
	db *sql.DB
}

func NewFileStore(db *sql.DB) FileService {
	return &fileStore{db: db}
}

func (f *fileStore) UploadFile(ctx context.Context, file *models.File) (*models.File, error) {
	query := `INSERT INTO files ("userID", name, path, size, "mimeType") VALUES ($1, $2, $3, $4, $5) RETURNING "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC"`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[UploadFile]: prepare statement", slog.Any("error", err))
		return nil, ErrFailedToUploadFile
	}
	defer stmt.Close()

	var newFile models.File
	if err := stmt.QueryRowContext(ctx, *file.UserID, *file.Name, *file.Path, *file.Size, *file.MimeType).Scan(
		&newFile.ID, &newFile.UserID, &newFile.Name, &newFile.Path, &newFile.Size, &newFile.MimeType, &newFile.CreatedAtUTC, &newFile.UpdatedAtUTC); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			switch pqErr.Constraint {
			case "files_user_name_unique":
				return nil, ErrFileNameAlreadyExists
			case "files_user_path_unique":
				return nil, ErrFilePathAlreadyExists
			default:
				slog.Error("[UploadFile]: unique constraint violation", slog.Any("error", err))
				return nil, ErrFailedToUploadFile
			}
		}
		slog.Error("[UploadFile]: query row", slog.Any("error", err))
		return nil, ErrFailedToUploadFile
	}

	return &newFile, nil
}

func (f *fileStore) GetFiles(ctx context.Context, userID string) ([]*models.ListFileResponse, error) {
	query := `SELECT "ID", name, size, "mimeType" FROM files WHERE "userID" = $1`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[GetFiles]: prepare statement", slog.Any("error", err))
		return nil, ErrFailedToGetFiles
	}
	defer stmt.Close()

	rows, err := stmt.QueryContext(ctx, userID)
	if err != nil {
		slog.Error("[GetFiles]: query rows", slog.Any("error", err))
		return nil, ErrFailedToGetFiles
	}
	defer rows.Close()

	var files []*models.ListFileResponse
	for rows.Next() {
		var file models.ListFileResponse
		if err := rows.Scan(&file.ID, &file.FileName, &file.FileSize, &file.MimeType); err != nil {
			slog.Error("[GetFiles]: scan row", slog.Any("error", err))
			return nil, ErrFailedToGetFiles
		}
		files = append(files, &file)
	}

	if err := rows.Err(); err != nil {
		slog.Error("[GetFiles]: rows error", slog.Any("error", err))
		return nil, ErrFailedToGetFiles
	}

	return files, nil
}

func (f *fileStore) GetFileByID(ctx context.Context, id string, userID string) (*models.File, error) {
	query := `SELECT "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC" FROM files WHERE "ID" = $1 AND "userID" = $2`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[GetFileByID]: prepare statement", slog.Any("error", err))
		return nil, ErrFailedToGetFileByID
	}
	defer stmt.Close()

	var file models.File
	if err := stmt.QueryRowContext(ctx, id, userID).Scan(
		&file.ID, &file.UserID, &file.Name, &file.Path, &file.Size, &file.MimeType, &file.CreatedAtUTC, &file.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrFileNotFound
		}
		slog.Error("[GetFileByID]: query row", slog.Any("error", err))
		return nil, ErrFailedToGetFileByID
	}

	return &file, nil
}

func (f *fileStore) UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) error {
	fields := make([]utils.UpdateField, 0, 4)
	if req.Name != nil {
		fields = append(fields, utils.UpdateField{Column: "name", Value: *req.Name})
	}
	if req.Path != nil {
		fields = append(fields, utils.UpdateField{Column: "path", Value: *req.Path})
	}

	query, args := utils.BuildUpdateSQL("files", fields, []string{"ID", "userID"})
	args[0] = id
	args[1] = userID

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[UpdateFile]: prepare statement", slog.Any("error", err))
		return ErrFailedToUpdateFile
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			switch pqErr.Constraint {
			case "files_user_name_unique":
				return ErrFileNameAlreadyExists
			case "files_user_path_unique":
				return ErrFilePathAlreadyExists
			default:
				slog.Error("[UpdateFile]: unique constraint violation", slog.Any("error", err))
				return ErrFailedToUpdateFile
			}
		}
		slog.Error("[UpdateFile]: execute statement", slog.Any("error", err))
		return ErrFailedToUpdateFile
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("[UpdateFile]: getting rows affected", slog.Any("error", err))
		return ErrFailedToUpdateFile
	}

	if rowsAffected == 0 {
		return ErrFileNotFound
	}

	return nil
}

func (f *fileStore) DeleteFile(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM files WHERE "ID" = $1 AND "userID" = $2`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[DeleteFile]: prepare statement", slog.Any("error", err))
		return ErrFailedToDeleteFile
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id, userID)
	if err != nil {
		slog.Error("[DeleteFile]: execute statement", slog.Any("error", err))
		return ErrFailedToDeleteFile
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("[DeleteFile]: getting rows affected", slog.Any("error", err))
		return ErrFailedToDeleteFile
	}

	if rowsAffected == 0 {
		return ErrFileNotFound
	}

	return nil
}
