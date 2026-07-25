package grpc

import (
	"context"
	"errors"
	"net/http"

	MMS "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	grpcUtils "github.com/rakshithrajs/cloud/UMS/internal/grpc/utils"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"google.golang.org/grpc/metadata"
)

func (c *Client) UploadFileGrpcHandler(ctx context.Context, userID, fileName string, content []byte) (*models.File, int, string) {
	ctx = metadata.AppendToOutgoingContext(ctx, config.UserIDMetadataKey, userID)

	resp, err := c.mmsClient.UploadFile(ctx, &MMS.UploadFileRequest{
		FileName: fileName,
		Content:  content,
	})
	if err != nil {
		status, msg := grpcUtils.MapGRPCError(err, handlerUtils.ErrFailedToUploadFile.Error())
		return nil, status, msg
	}

	fileID := resp.GetFile().GetID()
	if err := c.storage.CreateUserFile(ctx, userID, fileID, resp.GetFile().GetFileName()); err != nil {
		if errors.Is(err, handlerUtils.ErrUserFileAlreadyExists) {
			return nil, http.StatusConflict, err.Error()
		}
		if _, delErr := c.mmsClient.DeleteFile(ctx, &MMS.DeleteFileRequest{FileID: fileID}); delErr != nil {
			return nil, http.StatusInternalServerError, handlerUtils.ErrFailedToRollback.Error()
		}
		return nil, http.StatusInternalServerError, handlerUtils.ErrFailedToUploadFile.Error()
	}

	return &models.File{
		ID:       resp.GetFile().GetID(),
		FileName: resp.GetFile().GetFileName(),
		FileSize: resp.GetFile().GetFileSize(),
		MimeType: handlerUtils.MimeTypeToString(resp.GetFile().GetMimeType()),
	}, http.StatusCreated, config.NullString
}
