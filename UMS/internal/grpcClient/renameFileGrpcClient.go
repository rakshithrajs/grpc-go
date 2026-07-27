package grpcClient

import (
	"context"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *Client) RenameFileGrpcClient(ctx context.Context, userID, fileID, newName string) error {
	oldName, err := c.storage.UpdateUserFile(ctx, userID, fileID, newName)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	renameBody := &MMS.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	}
	if _, err := c.mmsClient.RenameFile(ctx, renameBody); err != nil {
		if _, rbErr := c.storage.UpdateUserFile(ctx, userID, fileID, oldName); rbErr != nil {
			return status.Error(codes.Internal, handlerErrors.ErrFailedToRollback.Error())
		}
		return err
	}

	return nil
}
