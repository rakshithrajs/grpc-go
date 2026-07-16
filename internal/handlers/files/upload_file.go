package handlers

import (
	filespb "cloud/gen/files/v1"
	"cloud/internal/config"
	"cloud/internal/handlers"
	"cloud/internal/models"
	"cloud/internal/storage"
	"cloud/internal/utils"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (f *FileHandler) UploadFile(ctx context.Context, req *filespb.UploadFileRequest) (*filespb.UploadFileResponse, error) {
	userID, err := handlers.UserIDFromContext(ctx)
	if err != nil {
		slog.Error("[UploadFile]: Missing user context", slog.Any("error", err))
		return nil, err
	}

	cleanedName := strings.TrimSpace(req.FileName)
	payload := models.UploadFileRequest{
		Name:     &cleanedName,
		Contents: req.Content,
	}

	if err := utils.Validate.Struct(&payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, utils.FirstError(err).Error())
	}

	config, err := config.GetConfig()
	if err != nil {
		slog.Error("[UploadFile]: Failed to get config", slog.Any("error", err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	fileSize := int64(len(payload.Contents))
	mimeType := http.DetectContentType(payload.Contents)
	userDir := filepath.Join(config.UserStoragePath, userID)
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		slog.Error("[UploadFile]: Failed to create user directory", slog.Any("error", err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}

	filePath := filepath.Join(userDir, cleanedName)

	fi, err := os.OpenFile(filePath, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0o644)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, status.Error(codes.AlreadyExists, storage.ErrFilePathAlreadyExists.Error())
		}
		slog.Error("[UploadFile]: Failed to create file", slog.Any("error", err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}
	defer fi.Close()

	size, err := fi.Write(payload.Contents)
	if err != nil || size != len(payload.Contents) {
		slog.Error("[UploadFile]: Failed to write file", slog.Any("error", err), slog.Int("size", size), slog.Int("expected", len(payload.Contents)))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
	}
	fi.Sync()

	file := &models.File{
		UserID:   &userID,
		Name:     &cleanedName,
		Path:     &filePath,
		Size:     &fileSize,
		MimeType: &mimeType,
	}

	savedFile, err := f.storage.UploadFile(ctx, file)
	if err != nil {
		if errors.Is(err, storage.ErrFailedToUploadFile) {
			return nil, status.Error(codes.Internal, storage.ErrFailedToUploadFile.Error())
		} else {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
	}
	return &filespb.UploadFileResponse{
		File: &filespb.File{
			ID:       *savedFile.ID,
			FileName: *savedFile.Name,
			FileSize: *savedFile.Size,
			MimeType: *savedFile.MimeType,
		},
	}, nil
}
