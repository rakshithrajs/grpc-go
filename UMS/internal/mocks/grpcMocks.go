package mocks

import (
	"context"

	MMSpb "github.com/rakshithrajs/cloud/UMS/gen/MMS/v1"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MockMMSClient struct {
	UploadGrpcErr   GrpcOperationError
	DownloadGrpcErr GrpcOperationError
	DeleteGrpcErr   GrpcOperationError
	RenameGrpcErr   GrpcOperationError
}

func (m *MockMMSClient) UploadFile(ctx context.Context, in *MMSpb.UploadFileRequest, opts ...grpc.CallOption) (*MMSpb.UploadFileResponse, error) {
	err := m.UploadGrpcErr

	switch err {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToUploadFile.Error())
	case GrpcOpFileAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, ErrFileAlreadyExists.Error())
	case GrpcOpRollbackFailure:
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToRollback.Error())
	}

	return &MMSpb.UploadFileResponse{
		File: &MMSpb.File{
			ID:       "550e8400-e29b-41d4-a716-446655440000",
			FileName: in.GetFileName(),
			FileSize: int64(len(in.GetContent())),
			MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
		},
	}, nil
}

func (m *MockMMSClient) DownloadFile(ctx context.Context, in *MMSpb.DownloadFileRequest, opts ...grpc.CallOption) (*MMSpb.DownloadFileResponse, error) {
	err := m.DownloadGrpcErr

	switch err {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToDownloadFile.Error())
	case GrpcOpNotFound:
		return &MMSpb.DownloadFileResponse{}, nil
	}

	return &MMSpb.DownloadFileResponse{
		FileName: "test-file.txt",
		MimeType: MMSpb.MimeType_MIME_TYPE_TEXT_PLAIN,
		Content:  []byte("test content"),
	}, nil
}

func (m *MockMMSClient) DeleteFile(ctx context.Context, in *MMSpb.DeleteFileRequest, opts ...grpc.CallOption) (*MMSpb.EmptyMessage, error) {
	err := m.DeleteGrpcErr

	switch err {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToDeleteFile.Error())
	case GrpcOpRollbackFailure:
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToRollback.Error())
	}

	return &MMSpb.EmptyMessage{}, nil
}

func (m *MockMMSClient) RenameFile(ctx context.Context, in *MMSpb.RenameFileRequest, opts ...grpc.CallOption) (*MMSpb.EmptyMessage, error) {
	err := m.RenameGrpcErr

	switch err {
	case GrpcOpMissingMetadata:
		return nil, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	case GrpcOpMissingUserID:
		return nil, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	case GrpcOpInternalError:
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToRenameFile.Error())
	case GrpcOpFileAlreadyExists:
		return nil, status.Error(codes.AlreadyExists, ErrFileAlreadyExists.Error())
	case GrpcOpNotFound:
		return &MMSpb.EmptyMessage{}, nil
	case GrpcOpRollbackFailure:
		return nil, status.Error(codes.Internal, handlerErrors.ErrFailedToRollback.Error())
	}

	return &MMSpb.EmptyMessage{}, nil
}
