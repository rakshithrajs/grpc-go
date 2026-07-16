package storage

import (
	"cloud/internal/models"
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/lib/pq"
	"github.com/lib/pq/pqerror"
)

var (
	ErrFailedToUploadFile    = errors.New("failed to upload file")
	ErrFileNameAlreadyExists = errors.New("file name already exists")
	ErrFilePathAlreadyExists = errors.New("file path already exists")
	ErrFailedToGetFiles      = errors.New("failed to get files")
	ErrFailedToGetFileByID   = errors.New("failed to get file by ID")
	ErrFileIDDoesntExist     = errors.New("file id doesn't exist")
	ErrFailedToUpdateFile    = errors.New("failed to update file")
	ErrFailedToDeleteFile    = errors.New("failed to delete file")
)

type FileSvc struct {
	db *sql.DB
}

func NewFileService(db *sql.DB) FileService {
	return &FileSvc{db: db}
}

func (f *FileSvc) UploadFile(ctx context.Context, file *models.File) (*models.File, error) {
	query := `INSERT INTO files ("userID", name, path, size, "mimeType") VALUES ($1, $2, $3, $4, $5) RETURNING "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC"`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[UploadFile]: prepare statement", slog.Any("error", err))
		return nil, ErrFailedToUploadFile
	}
	defer stmt.Close()

	var newFile models.File
	if err := stmt.QueryRowContext(ctx, file.UserID, file.Name, file.Path, file.Size, file.MimeType).Scan(&newFile.ID, &newFile.UserID, &newFile.Name, &newFile.Path, &newFile.Size, &newFile.MimeType, &newFile.CreatedAtUTC, &newFile.UpdatedAtUTC); err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == pqerror.UniqueViolation {
			switch pqErr.Constraint {
			case "files_user_name_unique":
				slog.Error("[UploadFile]: file name already exists for user", slog.Any("error", err))
				return nil, ErrFileNameAlreadyExists
			case "files_user_path_unique":
				slog.Error("[UploadFile]: file path already exists for user", slog.Any("error", err))
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

func (f *FileSvc) GetFiles(ctx context.Context, userID string) ([]*models.ListFileResponse, error) {
	query := `SELECT "ID", name,  size, "mimeType" FROM files WHERE "userID" = $1`

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

func (f *FileSvc) GetFileByID(ctx context.Context, id string, userID string) (*models.File, error) {
	query := `SELECT "ID", "userID", name, path, size, "mimeType", "createdAtUTC", "updatedAtUTC" FROM files WHERE "ID" = $1 AND "userID" = $2`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[GetFileByID]: prepare statement", slog.Any("error", err))
		return nil, ErrFailedToGetFileByID
	}

	var file models.File
	if err := stmt.QueryRowContext(ctx, id, userID).Scan(&file.ID, &file.UserID, &file.Name, &file.Path, &file.Size, &file.MimeType, &file.CreatedAtUTC, &file.UpdatedAtUTC); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("[GetFileByID]: No file with id found", slog.Any("error", err))
			return nil, ErrFileIDDoesntExist
		}
		slog.Error("[GetFileByID]: query row", slog.Any("error", err))
		return nil, ErrFailedToGetFileByID
	}

	return &file, nil
}

func (f *FileSvc) UpdateFile(ctx context.Context, id string, req models.UpdateFileRequest, userID string) error {
	fields := make([]UpdateField, 0, 4)
	if req.Name != nil {
		fields = append(fields, UpdateField{Column: "name", Value: req.Name})
	}
	if req.Path != nil {
		fields = append(fields, UpdateField{Column: "path", Value: req.Path})
	}

	query, args := BuildUpdateSQL("files", fields, []string{"ID", "userID"})
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
				slog.Error("[UpdateFile]: file name already exists for user", slog.Any("error", err))
				return ErrFileNameAlreadyExists
			case "files_user_path_unique":
				slog.Error("[UpdateFile]: file path already exists for user", slog.Any("error", err))
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
		slog.Error("[UpdateFile]: No file with id found")
	}

	return nil
}

func (f *FileSvc) DeleteFile(ctx context.Context, id string, userID string) error {
	query := `DELETE FROM files WHERE "ID" = $1 AND "userID" = $2`

	stmt, err := f.db.PrepareContext(ctx, query)
	if err != nil {
		slog.Error("[DeleteFile]: prepare statement", slog.Any("error", err))
		return ErrFailedToDeleteFile
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, id, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			slog.Error("[DeleteFile]: No file with id found", slog.Any("error", err))
			return ErrFileIDDoesntExist
		}
		slog.Error("[DeleteFile]: execute statement", slog.Any("error", err))
		return ErrFailedToDeleteFile
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		slog.Error("[DeleteFile]: getting rows affected", slog.Any("error", err))
		return ErrFailedToDeleteFile
	}

	if rowsAffected == 0 {
		slog.Error("[DeleteFile]: No file with id found", slog.Any("error", err))
	}

	return nil
}
