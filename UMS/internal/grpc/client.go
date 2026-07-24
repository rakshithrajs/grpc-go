package grpc

import (
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

type Client struct {
	mmsClient MMSpb.FilesClient
	storage   storage.UserFilesService
}

func NewClient(mmsClient MMSpb.FilesClient, storage storage.UserFilesService) *Client {
	return &Client{
		mmsClient: mmsClient,
		storage:   storage,
	}
}
