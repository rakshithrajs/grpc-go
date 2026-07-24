package grpc

import (
	"context"

	"github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (c *Client) RenameFileGrpcHandler(ctx context.Context, userID, fileID, newName string) error {
	oldName, err := c.storage.GetUserFileName(ctx, userID, fileID)
	if err != nil {
		return utils.ErrSomethingWentWrong
	}
	if oldName == "" {
		return nil
	}

	if err := c.storage.UpdateUserFile(ctx, userID, fileID, newName); err != nil {
		return utils.ErrSomethingWentWrong
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "userID", userID)

	if _, err := c.mmsClient.RenameFile(ctx, &MMS.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	}); err != nil {
		if rbErr := c.storage.UpdateUserFile(ctx, userID, fileID, oldName); rbErr != nil {
			_ = rbErr
		}
		return err
	}

	return nil
}
