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

func (c *Client) DeleteFileGrpcHandler(ctx context.Context, userID, fileID string) error {
	fileName, err := c.storage.DeleteUserFile(ctx, userID, fileID)
	if err != nil {
		return status.Error(codes.Internal, utils.ErrFailedToDeleteUserFile.Error())
	}
	if fileName == "" {
		return nil
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	if _, err := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); err != nil {
		if rbErr := c.storage.CreateUserFile(ctx, userID, fileID, fileName); rbErr != nil {
			return status.Error(codes.Internal, utils.ErrFailedToRollback.Error())
		}
		return err
	}

	return nil
}
