package grpc

import (
	"context"
	"strings"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *Client) DeleteFileGrpcHandler(ctx context.Context, userID, fileID string) error {
	if strings.TrimSpace(fileID) == config.NullString {
		return status.Error(codes.InvalidArgument, handlerUtils.ErrFileIDRequired.Error())
	}

	fileName, err := c.storage.DeleteUserFile(ctx, userID, fileID)
	if err != nil {
		return status.Error(codes.Internal, handlerUtils.ErrFailedToDeleteUserFile.Error())
	}
	if fileName == config.NullString {
		return nil
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	if _, err := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); err != nil {
		if rbErr := c.storage.CreateUserFile(ctx, userID, fileID, fileName); rbErr != nil {
			return status.Error(codes.Internal, handlerUtils.ErrFailedToRollback.Error())
		}
		return err
	}

	return nil
}
