package grpcClient

import (
	"context"
	"errors"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *Client) UploadFileGrpcClient(ctx context.Context, userID, fileName string, content []byte) (*models.File, error) {
	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	resp, err := c.mmsClient.UploadFile(ctx, &MMS.UploadFileRequest{
		FileName: fileName,
		Content:  content,
	})
	if err != nil {
		return nil, err
	}

	fileID := resp.GetFile().GetID()
	if err := c.storage.CreateUserFile(ctx, userID, fileID, resp.GetFile().GetFileName()); err != nil {
		if errors.Is(err, handlerErrors.ErrUserFileAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		if _, delErr := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); delErr != nil {
			return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToRollback.Error())
		}
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToUploadFile.Error())
	}

	return &models.File{
		ID:       resp.GetFile().GetID(),
		FileName: resp.GetFile().GetFileName(),
		FileSize: resp.GetFile().GetFileSize(),
		MimeType: handlerUtils.MimeTypeToString(resp.GetFile().GetMimeType()),
	}, nil
}
