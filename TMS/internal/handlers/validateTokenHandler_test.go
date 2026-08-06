package handlers

import (
	"context"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	TMSpb "github.com/rakshithrajs/cloud/TMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/TMS/internal/config"
	jwtHelper "github.com/rakshithrajs/cloud/TMS/internal/handlers/helpers"
	"google.golang.org/grpc/codes"
)

func TestValidateToken(t *testing.T) {
	config.SetConfig(&config.Config{JWTSecret: "test-secret"})

	tests := []struct {
		name         string
		token        string
		expectedCode codes.Code
		expectedErr  string
	}{
		{
			name:         "missing bearer prefix",
			token:        "token",
			expectedCode: codes.Unauthenticated,
			expectedErr:  jwtHelper.ErrMissingBearerToken.Error(),
		},
		{
			name:         "empty bearer token",
			token:        "Bearer ",
			expectedCode: codes.Unauthenticated,
			expectedErr:  jwtHelper.ErrInvalidToken.Error(),
		},
		{
			name:         "invalid token",
			token:        "Bearer invalid",
			expectedCode: codes.Unauthenticated,
			expectedErr:  jwtHelper.ErrInvalidToken.Error(),
		},
		{
			name:         "expired token",
			token:        "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "sub": "user-123", "iat": time.Now().Add(-time.Hour * 48).Unix(), "exp": time.Now().Add(-time.Hour * 24).Unix()}),
			expectedCode: codes.Unauthenticated,
			expectedErr:  jwtHelper.ErrTokenExpired.Error(),
		},
		{
			name:         "wrong issuer",
			token:        "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "wrong", "sub": "user-123", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix()}),
			expectedCode: codes.Unauthenticated,
			expectedErr:  jwtHelper.ErrInvalidToken.Error(),
		},
		{
			name:         "missing subject",
			token:        "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix()}),
			expectedCode: codes.Unauthenticated,
			expectedErr:  jwtHelper.ErrInvalidToken.Error(),
		},
		{
			name:         "future iat",
			token:        "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "sub": "user-123", "iat": time.Now().Add(time.Hour).Unix(), "exp": time.Now().Add(time.Hour * 2).Unix()}),
			expectedCode: codes.Unauthenticated,
			expectedErr:  jwtHelper.ErrInvalidToken.Error(),
		},
		{
			name:         "valid token",
			token:        "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "sub": "user-123", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix()}),
			expectedCode: codes.OK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTokenHandler()

			resp, err := handler.ValidateToken(context.Background(), &TMSpb.ValidateTokenRequest{AccessToken: tt.token})

			jwtHelper.CheckGRPCStatus(t, err, tt.expectedCode, tt.expectedErr)

			if tt.expectedCode == codes.OK {
				if resp.GetClaims().GetUserID() != "user-123" {
					t.Errorf("expected user ID 'user-123', got '%s'", resp.GetClaims().GetUserID())
				}
			}
		})
	}
}

func makeToken(secret string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}
