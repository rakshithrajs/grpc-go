package grpcClient

import (
	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	TMSpb "github.com/rakshithrajs/cloud/UMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

// MMSClient is a gRPC client wrapper for the Media Management Service.
type MMSClient struct {
	mmsClient MMSpb.FilesClient
	storage   storage.UserFilesService
}

// NewMMSClient creates a new MMSClient.
func NewMMSClient(mmsClient MMSpb.FilesClient, storage storage.UserFilesService) *MMSClient {
	return &MMSClient{
		mmsClient: mmsClient,
		storage:   storage,
	}
}

// TMSClient is a gRPC client wrapper for the Token Management Service.
type TMSClient struct {
	client TMSpb.TokensClient
}

// NewTMSClient creates a new TMSClient.
func NewTMSClient(client TMSpb.TokensClient) *TMSClient {
	return &TMSClient{client: client}
}
