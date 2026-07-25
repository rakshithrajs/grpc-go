package mocks

import (
	"context"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	modelUtils "github.com/rakshithrajs/cloud/UMS/internal/models/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockMMSClient struct {
	MockErr           GrpcOperationError
	MockDeleteErr     GrpcOperationError
	UploadReturnEmpty bool
}

func (m *MockMMSClient) UploadFile(ctx context.Context, in *MMSpb.UploadFileRequest, opts ...grpc.CallOption) (*MMSpb.UploadFileResponse, error) {
	switch m.MockErr {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerUtils.ErrFailedToUploadFile.Error())
	case GrpcOpFileNameAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, ErrFileNameAlreadyExists.Error())
	case GrpcOpFilePathAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, ErrFilePathAlreadyExists.Error())
	}

	if m.UploadReturnEmpty {
		return &MMSpb.UploadFileResponse{}, nil
	}

	return &MMSpb.UploadFileResponse{
		File: &MMSpb.File{
			ID:       "file-id-123",
			FileName: in.GetFileName(),
			FileSize: int64(len(in.GetContent())),
			MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
		},
	}, nil
}

func (m *MockMMSClient) DownloadFile(ctx context.Context, in *MMSpb.DownloadFileRequest, opts ...grpc.CallOption) (*MMSpb.DownloadFileResponse, error) {
	switch m.MockErr {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerUtils.ErrFailedToDownloadFile.Error())
	case GrpcOpNotFound:
		return &MMSpb.DownloadFileResponse{}, nil
	}

	return &MMSpb.DownloadFileResponse{
		FileName: "test-file.txt",
		MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
		Content:  []byte("test content"),
	}, nil
}

func (m *MockMMSClient) ListFiles(ctx context.Context, in *MMSpb.EmptyMessage, opts ...grpc.CallOption) (*MMSpb.ListFilesResponse, error) {
	switch m.MockErr {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, ErrFailedToGetFiles.Error())
	}

	return &MMSpb.ListFilesResponse{
		File: []*MMSpb.File{
			{ID: "file-1", FileName: "file1.txt", FileSize: 100, MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN},
			{ID: "file-2", FileName: "file2.txt", FileSize: 200, MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN},
		},
	}, nil
}

func (m *MockMMSClient) DeleteFile(ctx context.Context, in *MMSpb.DeleteFileRequest, opts ...grpc.CallOption) (*MMSpb.EmptyMessage, error) {
	err := m.MockDeleteErr
	if err == GrpcOpSuccess {
		err = m.MockErr
	}

	switch err {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, modelUtils.ErrFileIDRequired.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerUtils.ErrFailedToDeleteFile.Error())
	case GrpcOpNotFound:
		return &MMSpb.EmptyMessage{}, nil
	case GrpcOpRollbackFailure:
		return nil, status.Error(codes.Internal, "rollback failed")
	}

	return &MMSpb.EmptyMessage{}, nil
}

func (m *MockMMSClient) RenameFile(ctx context.Context, in *MMSpb.RenameFileRequest, opts ...grpc.CallOption) (*MMSpb.EmptyMessage, error) {
	err := m.MockErr

	switch err {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInvalidArgument:
		return nil, status.Error(codes.InvalidArgument, modelUtils.ErrFileIDRequired.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerUtils.ErrFailedToRenameFile.Error())
	case GrpcOpFileNameAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, ErrFileNameAlreadyExists.Error())
	case GrpcOpNotFound:
		return &MMSpb.EmptyMessage{}, nil
	case GrpcOpRollbackFailure:
		return nil, status.Error(codes.Internal, "rollback failed")
	}

	return &MMSpb.EmptyMessage{}, nil
}
