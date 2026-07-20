package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	filespb "github.com/rakshithrajs/cloud/services/files/gen/files/v1"
	"github.com/rakshithrajs/cloud/services/files/internal/models"
	"github.com/rakshithrajs/cloud/services/files/internal/storage"
	"github.com/rakshithrajs/cloud/services/files/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedToRenameFile = errors.New("failed to rename file")
)

func (f *FileHandler) RenameFile(ctx context.Context, req *filespb.RenameFileRequest) (*filespb.RenameFileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFileID() == "" {
		return nil, status.Error(codes.InvalidArgument, ErrFileIDRequired.Error())
	}

	newName := strings.TrimSpace(req.GetNewName())
	payload := models.RenameFileRequest{Name: &newName}
	if err := utils.Validate.Struct(payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, strings.Join(utils.Errors(err), "; "))
	}

	file, err := f.fileService.GetFileByID(ctx, req.GetFileID(), userID)
	if err != nil {
		if errors.Is(err, storage.ErrFileNotFound) {
			return &filespb.RenameFileResponse{}, nil
		}
		slog.Error("[RenameFile]: failed to get file", slog.Any("error", err))
		return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
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

	if err := f.fileService.UpdateFile(ctx, req.GetFileID(), updateReq, userID); err != nil {
		if rbErr := os.Rename(newPath, oldPath); rbErr != nil {
			slog.Error("[RenameFile]: failed to rollback disk rename", slog.Any("error", rbErr), slog.String("oldPath", oldPath), slog.String("newPath", newPath))
		}

		switch {
		case errors.Is(err, storage.ErrFileNotFound):
			return &filespb.RenameFileResponse{}, nil
		case errors.Is(err, storage.ErrFileNameAlreadyExists), errors.Is(err, storage.ErrFilePathAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			slog.Error("[RenameFile]: failed to update file record", slog.Any("error", err))
			return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
		}
	}

	return &filespb.RenameFileResponse{}, nil
}
