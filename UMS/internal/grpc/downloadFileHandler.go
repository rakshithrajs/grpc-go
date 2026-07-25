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

func (c *Client) DownloadFileGrpcHandler(ctx context.Context, userID, fileID string) (*MMS.DownloadFileResponse, error) {
	if strings.TrimSpace(fileID) == "" {
		return nil, status.Error(codes.InvalidArgument, handlerUtils.ErrFileIDRequired.Error())
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	resp, err := c.mmsClient.DownloadFile(ctx, &MMS.DownloadFileRequest{FileID: fileID})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
