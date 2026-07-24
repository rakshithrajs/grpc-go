package grpc

import (
	"context"

	"github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (c *Client) DeleteFileGrpcHandler(ctx context.Context, userID, fileID string) error {
	fileName, err := c.storage.DeleteUserFile(ctx, userID, fileID)
	if err != nil {
		return utils.ErrSomethingWentWrong
	}
	if fileName == "" {
		return nil
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "userID", userID)

	if _, err := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); err != nil {
		if rbErr := c.storage.CreateUserFile(ctx, userID, fileID, fileName); rbErr != nil {
			return rbErr
		}
		return err
	}

	return nil
}
