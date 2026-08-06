package grpcClient

import (
	"context"
	"testing"

	TMSpb "github.com/rakshithrajs/cloud/UMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	handlerErrors "github.com/rakshithrajs/cloud/UMS/internal/handlers/errors"
	"github.com/rakshithrajs/cloud/UMS/internal/mocks"
	mockUtils "github.com/rakshithrajs/cloud/UMS/internal/mocks/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		grpcErr      mocks.GrpcOperationError
		expectedCode codes.Code
		expectedErr  string
		expectedData string
	}{
		{
			name:         "generate token returns access token",
			userID:       "test-user-id",
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
			expectedData: "test-access-token",
		},
		{
			name:         "generate token fails when TMS returns internal error",
			userID:       "test-user-id",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  mocks.ErrFailedToGenerateToken.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokensClient := &mocks.MockTokensClient{
				AccessToken: "test-access-token",
				GenerateErr: tt.grpcErr,
			}
			c := NewTMSClient(tokensClient)

			got, err := c.GenerateToken(context.Background(), tt.userID)

			st, _ := status.FromError(err)

			mockUtils.CheckData(t, st.Code(), tt.expectedCode)
			mockUtils.CheckError(t, st.Message(), tt.expectedErr)
			mockUtils.CheckData(t, tokensClient.UserID, tt.userID)

			if st.Code() == codes.OK {
				mockUtils.CheckData(t, got, tt.expectedData)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	tests := []struct {
		name         string
		accessToken  string
		grpcErr      mocks.GrpcOperationError
		expectedCode codes.Code
		expectedErr  string
		expectedData *TokenClaims
	}{
		{
			name:         "validate token returns claims",
			accessToken:  "Bearer valid-token",
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
			expectedData: &TokenClaims{
				UserID:   "test-user-id",
				Issuer:   "cloud-app",
				IssuedAt: 1234567890,
			},
		},
		{
			name:         "validate token fails when TMS returns missing bearer token error",
			accessToken:  "invalid-token",
			grpcErr:      mocks.GrpcOpMissingBearerToken,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrMissingBearerToken.Error(),
		},
		{
			name:         "validate token fails when TMS returns missing bearer token error",
			accessToken:  "Bearer ",
			grpcErr:      mocks.GrpcOpInvalidToken,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrInvalidToken.Error(),
		},
		{
			name:         "validate token fails when TMS returns invalid token error",
			accessToken:  "Bearer invalid-token",
			grpcErr:      mocks.GrpcOpInvalidToken,
			expectedCode: codes.Unauthenticated,
			expectedErr:  mocks.ErrInvalidToken.Error(),
		},
		{
			name:         "validate token fails when TMS returns internal error",
			accessToken:  "Bearer valid-token",
			grpcErr:      mocks.GrpcOpInternalError,
			expectedCode: codes.Internal,
			expectedErr:  handlerErrors.ErrSomethingWentWrong.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokensClient := &mocks.MockTokensClient{
				Claims: &TMSpb.TokenClaims{
					UserID:   "test-user-id",
					Issuer:   "cloud-app",
					IssuedAt: 1234567890,
				},
				ValidateErr: tt.grpcErr,
			}
			c := NewTMSClient(tokensClient)

			got, err := c.ValidateToken(context.Background(), tt.accessToken)

			st, _ := status.FromError(err)

			mockUtils.CheckData(t, st.Code(), tt.expectedCode)
			mockUtils.CheckError(t, st.Message(), tt.expectedErr)
			mockUtils.CheckData(t, tokensClient.AccessTokenRequest, tt.accessToken)

			if st.Code() == codes.OK {
				mockUtils.CheckData(t, got, tt.expectedData)
			}
		})
	}
}
