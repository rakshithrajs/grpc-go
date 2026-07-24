package grpc

import (
	"context"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"google.golang.org/grpc/metadata"
)

func (c *Client) DownloadFileGrpcHandler(ctx context.Context, userID, fileID string) (*MMS.DownloadFileResponse, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	resp, err := c.mmsClient.DownloadFile(ctx, &MMS.DownloadFileRequest{FileID: fileID})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
