package grpc

import (
	"context"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"google.golang.org/grpc/metadata"
)

func (c *Client) DownloadFileGrpcHandler(ctx context.Context, userID, fileID string) (string, MMS.MimeType, []byte, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, "userID", userID)

	resp, err := c.mmsClient.DownloadFile(ctx, &MMS.DownloadFileRequest{FileID: fileID})
	if err != nil {
		return "", MMS.MimeType_MIME_TYPE_UNSPECIFIED, nil, err
	}

	return resp.GetFileName(), resp.GetMimeType(), resp.GetContent(), nil
}
