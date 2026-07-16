package handlers

import (
	filespb "cloud/gen/files/v1"
	"cloud/internal/handlers"
	"cloud/internal/models"
	"cloud/internal/storage"
	"cloud/internal/utils"
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedToRenameFile = errors.New("failed to rename file")
)

func (f *FileHandler) RenameFile(ctx context.Context, req *filespb.RenameFileRequest) (*filespb.RenameFileResponse, error) {
	userID, err := handlers.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFileID() == "" {
		return nil, status.Error(codes.InvalidArgument, ErrFileIDRequired.Error())
	}

	newName := strings.TrimSpace(req.GetNewName())
	payload := models.RenameFileRequest{Name: &newName}
	if err := utils.Validate.Struct(&payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, utils.FirstError(err).Error())
	}

	file, err := f.storage.GetFileByID(ctx, req.GetFileID(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileIDDoesntExist) {
			return &filespb.RenameFileResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	oldPath := *file.Path
	userDir := filepath.Dir(oldPath)
	newPath := filepath.Join(userDir, newName)

	if oldPath == newPath {
		return &filespb.RenameFileResponse{}, nil
	}

	if _, err := os.Stat(newPath); err == nil {
		return nil, status.Error(codes.AlreadyExists, storage.ErrFilePathAlreadyExists.Error())
	} else if !errors.Is(err, os.ErrNotExist) {
		slog.Error("[RenameFile]: failed to check target path", slog.Any("error", err), slog.String("path", newPath))
		return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		slog.Error("[RenameFile]: failed to rename file on disk", slog.Any("error", err), slog.String("oldPath", oldPath), slog.String("newPath", newPath))
		return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
	}

	updateReq := models.UpdateFileRequest{
		Name: &newName,
		Path: &newPath,
	}
	if err := f.storage.UpdateFile(ctx, req.GetFileID(), updateReq, userID); err != nil {
		if rbErr := os.Rename(newPath, oldPath); rbErr != nil {
			slog.Error("[RenameFile]: failed to rollback disk rename", slog.Any("error", rbErr), slog.String("oldPath", oldPath), slog.String("newPath", newPath))
		}

		switch {
		case errors.Is(err, storage.ErrFileIDDoesntExist):
			return &filespb.RenameFileResponse{}, nil
		case errors.Is(err, storage.ErrFileNameAlreadyExists), errors.Is(err, storage.ErrFilePathAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			slog.Error("[RenameFile]: failed to update file record", slog.Any("error", err), slog.String("fileID", req.GetFileID()))
			return nil, status.Error(codes.Internal, storage.ErrFailedToUpdateFile.Error())
		}
	}

	return &filespb.RenameFileResponse{}, nil
}
