package grpc

import (
	"context"
	"reflect"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestClient_UploadFileGrpcHandler(t *testing.T) {
	tests := []struct {
		name              string
		content           []byte
		grpcErr           mocks.GrpcOperationError
		createDbErr       mocks.DbOperationError
		uploadReturnEmpty bool
		expectedCode      codes.Code
		expectedErr       string
		expectedFile      *models.File
	}{
		{
			name:         "upload fails when content is empty",
			content:      []byte{},
			expectedCode: codes.InvalidArgument,
			expectedErr:  utils.ErrFileIsRequired.Error(),
		},
		{
			name:         "upload fails due to missing metadata",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingMetadata,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingMetadata.Error(),
		},
		{
			name:         "upload fails due to missing user id",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpMissingUserID,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingUserID.Error(),
		},
		{
			name:         "upload fails due to internal error",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails as file name already exists",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFileNameAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFileNameAlreadyExists.Error(),
		},
		{
			name:         "upload fails as file path already exists",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpFilePathAlreadyExists,
			expectedCode: codes.AlreadyExists,
			expectedErr:  mocks.ErrFilePathAlreadyExists.Error(),
		},
		{
			name:              "upload fails when grpc returns empty file",
			content:           []byte("test content"),
			uploadReturnEmpty: true,
			expectedCode:      codes.Internal,
			expectedErr:       utils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails when user file mapping already exists",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpDuplicateFile,
			expectedCode: codes.AlreadyExists,
			expectedErr:  utils.ErrUserFileAlreadyExists.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback succeeds",
			content:      []byte("test content"),
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  utils.ErrFailedToUploadFile.Error(),
		},
		{
			name:         "upload fails due to db internal error and rollback fails",
			content:      []byte("test content"),
			grpcErr:      mocks.GrpcOpRollbackFailure,
			createDbErr:  mocks.DbOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  utils.ErrFailedToRollback.Error(),
		},
		{
			name:         "upload succeeds",
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

			file, err := c.UploadFileGrpcHandler(context.Background(), "user-123", "test.txt", tt.content)

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

			if !reflect.DeepEqual(file, tt.expectedFile) {
				t.Errorf("expected file %+v got %+v", tt.expectedFile, file)
			}
		})
	}
}
