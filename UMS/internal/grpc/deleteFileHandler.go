package grpc

import (
	"context"
	"net/http"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpcUtils "github.com/rakshithrajs/cloud/UMS/internal/grpc/utils"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"google.golang.org/grpc/metadata"
)

func (c *Client) DeleteFileGrpcHandler(ctx context.Context, userID, fileID string) (int, string) {
	fileName, err := c.storage.DeleteUserFile(ctx, userID, fileID)
	if err != nil {
		return http.StatusBadRequest, err.Error()
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	if _, err := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); err != nil {
		if rbErr := c.storage.CreateUserFile(ctx, userID, fileID, fileName); rbErr != nil {
			return http.StatusInternalServerError, handlerErrors.ErrFailedToRollback.Error()
		}
		return grpcUtils.MapGRPCError(err, storage.ErrFailedToDeleteUserFile.Error())
	}

	return http.StatusOK, config.NullString
}
