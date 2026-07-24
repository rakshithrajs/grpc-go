package grpc

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/metadata"
)

func (c *Client) UploadFileGrpcHandler(ctx context.Context, userID, fileName string, content []byte) (*models.File, error) {
	req := models.UploadRequest{UserID: userID, FileName: fileName}
	if err := utils.Validate.Struct(req); err != nil {
		var verrs validator.ValidationErrors
		if errors.As(err, &verrs) {
			for _, e := range verrs {
				switch e.StructField() {
				case "UserID":
					return nil, utils.ErrUnauthorized
				case "FileName":
					return nil, utils.ErrFileIsRequired
				}
			}
		}
		return nil, utils.ErrSomethingWentWrong
	}
	if len(content) == 0 {
		return nil, utils.ErrFileIsRequired
	}

	ctx = metadata.AppendToOutgoingContext(ctx, "userID", userID)

	resp, err := c.mmsClient.UploadFile(ctx, &MMS.UploadFileRequest{
		FileName: fileName,
		Content:  content,
	})
	if err != nil {
		return nil, err
	}

	if resp.GetFile() == nil || resp.GetFile().GetID() == "" {
		return nil, utils.ErrFailedToUploadFile
	}

	fileID := resp.GetFile().GetID()
	if err := c.storage.CreateUserFile(ctx, userID, fileID, resp.GetFile().GetFileName()); err != nil {
		if _, delErr := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); delErr != nil {
			_ = delErr
		}
		return nil, utils.ErrFailedToUploadFile
	}

	return &models.File{
		ID:       resp.GetFile().GetID(),
		FileName: resp.GetFile().GetFileName(),
		FileSize: resp.GetFile().GetFileSize(),
		MimeType: utils.MimeTypeToString(resp.GetFile().GetMimeType()),
	}, nil
}
