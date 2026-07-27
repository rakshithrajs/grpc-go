package grpc

import (
	"context"
	"net/http"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpcUtils "github.com/rakshithrajs/cloud/UMS/internal/grpc/utils"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"google.golang.org/grpc/metadata"
)

func (c *Client) RenameFileGrpcHandler(ctx context.Context, userID, fileID, newName string) (int, string) {
	oldName, err := c.storage.UpdateUserFile(ctx, userID, fileID, newName)
	if err != nil {
		return http.StatusInternalServerError, err.Error()
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	renameBody := &MMS.RenameFileRequest{
		FileID:  fileID,
		NewName: newName,
	}
	if _, err := c.mmsClient.RenameFile(ctx, renameBody); err != nil {
		if _, rbErr := c.storage.UpdateUserFile(ctx, userID, fileID, oldName); rbErr != nil {
			return http.StatusInternalServerError, handlerErrors.ErrFailedToRollback.Error()
		}
		return grpcUtils.MapGRPCError(err, handlerErrors.ErrFailedToRenameFile.Error())
	}

	return http.StatusOK, config.NullString
}
