package grpc

import (
	"context"
	"strings"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *Client) RenameFileGrpcHandler(ctx context.Context, userID, fileID, newName string) error {
	if strings.TrimSpace(fileID) == config.NullString {
		return status.Error(codes.InvalidArgument, handlerUtils.ErrFileIDRequired.Error())
	}
	if strings.TrimSpace(newName) == config.NullString {
		return status.Error(codes.InvalidArgument, modelUtils.ErrNewNameRequired.Error())
	}

	oldName, err := c.storage.UpdateUserFile(ctx, userID, fileID, newName)
	if err != nil {
		return status.Error(codes.Internal, handlerUtils.ErrFailedToUpdateUserFile.Error())
	}
	if oldName == config.NullString {
		return nil
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	if _, err := c.mmsClient.RenameFile(ctx, &MMS.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	}); err != nil {
		if _, rbErr := c.storage.UpdateUserFile(ctx, userID, fileID, oldName); rbErr != nil {
			return status.Error(codes.Internal, handlerUtils.ErrFailedToRollback.Error())
		}
		return err
	}

	return nil
}
