package grpc

import (
	"context"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerUtils "github.com/rakshithrajs/cloud/UMS/internal/handlers/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClient_UploadFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		fileName          string
		content           []byte
		grpcErr           mocks.GrpcOperationError
		createDbErr       mocks.DbOperationError
		uploadReturnEmpty bool
		expectedCode      codes.Code
		expectedErr       string
		expectedFile      *models.File
	}{
		{
			name:         "upload fails when file name is empty",
			content:      []byte("test content"),
			fileName:     config.NullString,
			expectedCode: codes.InvalidArgument,
			expectedErr:  handlerUtils.ErrFileNameRequired.Error(),
		},
		{
			name:         "upload fails when file name is whitespace",
			content:      []byte("test content"),
			fileName:     "   ",
			expectedCode: codes.InvalidArgument,
			expectedErr:  handlerUtils.ErrFileNameRequired.Error(),
		},
		{
			name:         "upload fails when content is empty",
			fileName:     "test.txt",
			content:      []byte{},
			expectedCode: codes.InvalidArgument,
			expectedErr:  handlerUtils.ErrFileIsRequired.Error(),
		},
		{
			name:         "upload fails due to missing metadata",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "upload fails due to missing user id",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "upload fails due to internal error",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails as file name already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFileNameAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "upload fails as file path already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFilePathAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFilePathAlreadyExists.Error(),
		},
		{
			name:              "upload fails when grpc returns empty file",
			fileName:          "test.txt",
			content:           []byte("test content"),
			uploadReturnEmpty: true,
			expectedCode:      codes.Internal,
			expectedErr:       handlerUtils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails when user file mapping already exists",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpDuplicateFile,
			expectedCode: codes.AlreadyExists,
			expectedErr:  handlerUtils.ErrUserFileAlreadyExists.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback succeeds",
			fileName:     "test.txt",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerUtils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback fails",
			fileName:     "test.txt",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpRollbackFailure,
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerUtils.ErrFailedToRollback.Error(),
		},
		{
			name:         "upload succeeds",
			fileName:     "test.txt",
			content:      []byte("test content"),
			expectedCode: codes.OK,
			expectedFile: &models.File{
				ID:       "file-id-123",
				FileName: "test.txt",
				FileSize: 12,
				MimeType: "text/plain",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mmsClient := &mocks.MockMMSClient{
				MockErr:           tt.grpcErr,
				UploadReturnEmpty: tt.uploadReturnEmpty,
			}
			svc := &mocks.MockUserFilesService{CreateUserFileError: tt.createDbErr}
			c := NewClient(mmsClient, svc)

			file, err := c.UploadFileGrpcHandler(context.Background(), "user-123", tt.fileName, tt.content)

			st, _ := status.FromError(err)

			if tt.expectedCode != st.Code() {
				t.Errorf("expected code %v got %v", tt.expectedCode, st.Code())
			}

			if st.Code() != codes.OK {
				if tt.expectedErr != st.Message() {
					t.Errorf("expected error %v got %v", tt.expectedErr, st.Message())
				}
				return
			}

			mockUtils.CheckData(t, file, tt.expectedFile)
		})
	}
}
