package handlers

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) UploadFile(ctx context.Context, req *MMSpb.UploadFileRequest) (*MMSpb.UploadFileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	cleanedName := strings.TrimSpace(req.GetFileName())

	payload := models.UploadFileRequest{
		Name:     cleanedName,
		Contents: req.GetContent(),
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(logPrefix(fnUploadFile)+"failed to get config", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	fileSize := int64(len(payload.Contents))
	mimeType := models.ParseMimeType(http.DetectContentType(payload.Contents))

	userDir := filepath.Join(cfg.UserStoragePath, userID)
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		slog.Error(logPrefix(fnUploadFile)+"failed to create user directory", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	filePath := filepath.Join(userDir, cleanedName)

	fi, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, status.Error(codes.AlreadyExists, storage.ErrFileAlreadyExists.Error())
		}
		slog.Error(logPrefix(fnUploadFile)+"failed to create file", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}
	defer fi.Close()

	if written, err := fi.Write(payload.Contents); err != nil {
		slog.Error(logPrefix(fnUploadFile)+"failed to write file", slog.Any(config.ErrorKey, err), slog.Int("written", written), slog.Int("expected", len(payload.Contents)))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	if err := fi.Sync(); err != nil {
		slog.Error(logPrefix(fnUploadFile)+"failed to sync file", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	if err := fi.Close(); err != nil {
		slog.Error(logPrefix(fnUploadFile)+"failed to close file", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	file := &models.File{
		UserID:   userID,
		Name:     cleanedName,
		Path:     filePath,
		Size:     fileSize,
		MimeType: mimeType,
	}

	savedFile, err := f.fileService.UploadFile(ctx, file)
	if err != nil {
		if errors.Is(err, storage.ErrFileAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if err := os.Remove(filePath); err != nil {
			slog.Error(logPrefix(fnUploadFile)+"failed to remove file", slog.Any(config.ErrorKey, err))
			return nil, status.Error(codes.Internal, ErrFailedToRollback.Error())
		}
		slog.Error(logPrefix(fnUploadFile)+"failed to save file metadata", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	return &MMSpb.UploadFileResponse{
		File: &MMSpb.File{
			ID:       savedFile.ID,
			FileName: savedFile.Name,
			FileSize: savedFile.Size,
			MimeType: toProtoMimeType(savedFile.MimeType),
		},
	}, nil
}
