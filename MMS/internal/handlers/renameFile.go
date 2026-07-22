package handlers

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	MMSpb "github.com/rakshithrajs/cloud/MMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/MMS/internal/models"
	"github.com/rakshithrajs/cloud/MMS/internal/storage"
	"github.com/rakshithrajs/cloud/MMS/internal/utils"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedToRenameFile = errors.New("failed to rename file")
)

func (f *FileHandler) RenameFile(ctx context.Context, req *MMSpb.RenameFileRequest) (*MMSpb.RenameFileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	if req.GetFileID() == nullString {
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
			return &MMSpb.RenameFileResponse{}, nil
		}
		slog.Error(logPrefix(fnRenameFile)+"failed to get file", slog.Any("error", err))
		return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
	}

	oldPath := *file.Path
	userDir := filepath.Dir(oldPath)
	newPath := filepath.Join(userDir, newName)

	if *file.Name == newName {
		return &MMSpb.RenameFileResponse{}, nil
	}

	if err := os.Rename(oldPath, newPath); err != nil {
		slog.Error(logPrefix(fnRenameFile)+"failed to rename file on disk", slog.Any("error", err), slog.String("oldPath", oldPath), slog.String("newPath", newPath))
		return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
	}

	if err := f.fileService.UpdateFile(ctx, req.GetFileID(), models.UpdateFileRequest{Name: &newName, Path: &newPath}, userID); err != nil {
		slog.Error(logPrefix(fnRenameFile)+"failed to update file record", slog.Any("error", err))
		if rbErr := os.Rename(newPath, oldPath); rbErr != nil {
			slog.Error(logPrefix(fnRenameFile)+"failed to rollback disk rename", slog.Any("error", rbErr), slog.String("oldPath", oldPath), slog.String("newPath", newPath))
		}
		switch {
		case errors.Is(err, storage.ErrFileNotFound):
			return &MMSpb.RenameFileResponse{}, nil
		case errors.Is(err, storage.ErrFileNameAlreadyExists):
			return nil, status.Error(codes.AlreadyExists, err.Error())
		default:
			return nil, status.Error(codes.Internal, ErrFailedToRenameFile.Error())
		}
	}

	return &MMSpb.RenameFileResponse{}, nil
}
