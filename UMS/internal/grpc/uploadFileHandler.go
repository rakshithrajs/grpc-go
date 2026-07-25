package grpc

import (
	"context"
	"errors"
	"strings"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *Client) UploadFileGrpcHandler(ctx context.Context, userID, fileName string, content []byte) (*models.File, error) {
	if strings.TrimSpace(fileName) == "" {
		return nil, status.Error(codes.InvalidArgument, handlerUtils.ErrFileNameRequired.Error())
	}
	if len(content) == 0 {
		return nil, status.Error(codes.InvalidArgument, handlerUtils.ErrFileIsRequired.Error())
	}

	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	resp, err := c.mmsClient.UploadFile(ctx, &MMS.UploadFileRequest{
		FileName: fileName,
		Content:  content,
	})
	if err != nil {
		return nil, err
	}

	if resp.GetFile() == nil || resp.GetFile().GetID() == "" {
		return nil, status.Error(codes.Internal, handlerUtils.ErrFailedToUploadFile.Error())
	}

	fileID := resp.GetFile().GetID()
	if err := c.storage.CreateUserFile(ctx, userID, fileID, resp.GetFile().GetFileName()); err != nil {
		if errors.Is(err, handlerUtils.ErrUserFileAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, handlerUtils.ErrUserFileAlreadyExists.Error())
		}
		if _, delErr := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); delErr != nil {
			return nil, status.Error(codes.Internal, handlerUtils.ErrFailedToRollback.Error())
		}
		return nil, status.Error(codes.Internal, handlerUtils.ErrFailedToUploadFile.Error())
	}

	return &models.File{
		ID:       resp.GetFile().GetID(),
		FileName: resp.GetFile().GetFileName(),
		FileSize: resp.GetFile().GetFileSize(),
		MimeType: handlerUtils.MimeTypeToString(resp.GetFile().GetMimeType()),
	}, nil
}
