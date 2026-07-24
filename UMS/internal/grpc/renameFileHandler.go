package grpc

import (
	"context"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *Client) RenameFileGrpcHandler(ctx context.Context, userID, fileID, newName string) error {
	oldName, err := c.storage.UpdateUserFile(ctx, userID, fileID, newName)
	if err != nil {
		return status.Error(codes.Internal, utils.ErrFailedToUpdateUserFile.Error())
	}
	if oldName == "" {
		return nil
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	if _, err := c.mmsClient.RenameFile(ctx, &MMS.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	}); err != nil {
		if _, rbErr := c.storage.UpdateUserFile(ctx, userID, fileID, oldName); rbErr != nil {
			return status.Error(codes.Internal, utils.ErrFailedToRollback.Error())
		}
		return err
	}

	return nil
}
