package grpc

import (
	"context"
	"errors"
	"strings"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func (c *Client) UploadFileGrpcHandler(ctx context.Context, userID, fileName string, content []byte) (*models.File, error) {
	if strings.TrimSpace(fileName) == "" {
		return nil, status.Error(codes.InvalidArgument, utils.ErrFileNameRequired.Error())
	}
	if len(content) == 0 {
		return nil, status.Error(codes.InvalidArgument, utils.ErrFileIsRequired.Error())
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
		return nil, status.Error(codes.Internal, utils.ErrFailedToUploadFile.Error())
	}

	fileID := resp.GetFile().GetID()
	if err := c.storage.CreateUserFile(ctx, userID, fileID, resp.GetFile().GetFileName()); err != nil {
		if errors.Is(err, utils.ErrUserFileAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, utils.ErrUserFileAlreadyExists.Error())
		}
		if _, delErr := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); delErr != nil {
			return nil, status.Error(codes.Internal, utils.ErrFailedToRollback.Error())
		}
		return nil, status.Error(codes.Internal, utils.ErrFailedToUploadFile.Error())
	}

	return &models.File{
		ID:       resp.GetFile().GetID(),
		FileName: resp.GetFile().GetFileName(),
		FileSize: resp.GetFile().GetFileSize(),
		MimeType: utils.MimeTypeToString(resp.GetFile().GetMimeType()),
	}, nil
}
