package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/config"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RenameFile renames a file in storage for the authenticated user.
func (f *FileHandler) RenameFile(ctx context.Context, req *MMSpb.RenameFileRequest) (*MMSpb.EmptyMessage, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	updateBody := models.UpdateFileRequest{
		Name: req.GetNewName(),
	}

	file, err := f.fileService.UpdateFile(ctx, req.GetFileID(), updateBody, userID)
	if err != nil {
		slog.Error(logPrefix(fnRenameFile)+"failed to update file record", slog.Any(config.ErrorKey, err))
		if errors.Is(err, storage.ErrFileAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
	}

	if file == nil || file.ID == config.NullString {
		return &MMSpb.EmptyMessage{}, nil
	}

	oldPath := file.Path
	userDir := filepath.Dir(oldPath)
	newPath := filepath.Join(userDir, updateBody.Name)
	updateBody.Path = newPath

	if oldPath == newPath {
		return &MMSpb.EmptyMessage{}, nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		slog.Error(logPrefix(fnRenameFile)+"failed to rename file on disk", slog.Any(config.ErrorKey, err), slog.String("oldPath", oldPath), slog.String("newPath", newPath))
		updateBody.Name = file.Name
		updateBody.Path = file.Path
		if _, rbErr := f.fileService.UpdateFile(ctx, req.GetFileID(), updateBody, userID); rbErr != nil {
			slog.Error(logPrefix(fnRenameFile)+"failed to roll back file rename", slog.Any(config.ErrorKey, rbErr), slog.String("oldPath", oldPath), slog.String("newPath", newPath))
			return nil, status.Error(codes.Internal, ErrFailedToRollback.Error())
		}
		return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
	}

	return &MMSpb.EmptyMessage{}, nil
}
