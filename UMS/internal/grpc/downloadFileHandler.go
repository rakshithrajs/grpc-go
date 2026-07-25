package grpc

import (
	"context"
	"net/http"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpcUtils "github.com/rakshithrajs/cloud/UMS/internal/grpc/utils"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"google.golang.org/grpc/metadata"
)

func (c *Client) DownloadFileGrpcHandler(ctx context.Context, userID, fileID string) (*MMSpb.DownloadFileResponse, int, string) {
	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	resp, err := c.mmsClient.DownloadFile(ctx, &MMS.DownloadFileRequest{FileID: fileID})
	if err != nil {
		status, msg := grpcUtils.MapGRPCError(err, handlerUtils.ErrFailedToDownloadFile.Error())
		return &MMSpb.DownloadFileResponse{}, status, msg
	}

	return resp, http.StatusOK, config.NullString
}
