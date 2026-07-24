package grpc

import (
	"context"

	"github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (c *Client) RenameFileGrpcHandler(ctx context.Context, userID, fileID, newName string) error {
	oldName, err := c.storage.UpdateUserFile(ctx, userID, fileID, newName)
	if err != nil {
		return utils.ErrSomethingWentWrong
	}
	if oldName == "" {
		return nil
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "userID", userID)

	if _, err := c.mmsClient.RenameFile(ctx, &MMS.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	}); err != nil {
		if _, rbErr := c.storage.UpdateUserFile(ctx, userID, fileID, oldName); rbErr != nil {
			return rbErr
		}
		return err
	}

	return nil
}
