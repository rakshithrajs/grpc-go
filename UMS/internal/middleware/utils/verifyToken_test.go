package middleware

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
)

func TestVerifyToken(t *testing.T) {
	config.SetConfig(&config.Config{JWTSecret: "test-secret"})

	tests := []struct {
		name    string
		token   string
		wantErr error
	}{
		{
			name:    "missing bearer prefix",
			token:   "token",
			wantErr: ErrMissingBearerToken,
		},
		{
			name:    "empty bearer token",
			token:   "Bearer ",
			wantErr: ErrMissingBearerToken,
		},
		{
			name:    "invalid token",
			token:   "Bearer invalid",
			wantErr: ErrInvalidToken,
		},
		{
			name:    "expired token",
			token:   "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "sub": "user-123", "iat": time.Now().Add(-time.Hour * 48).Unix(), "exp": time.Now().Add(-time.Hour * 24).Unix()}),
			wantErr: ErrTokenExpired,
		},
		{
			name:    "wrong issuer",
			token:   "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "wrong", "sub": "user-123", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix()}),
			wantErr: ErrInvalidToken,
		},
		{
			name:    "missing subject",
			token:   "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix()}),
			wantErr: ErrInvalidToken,
		},
		{
			name:    "future iat",
			token:   "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "sub": "user-123", "iat": time.Now().Add(time.Hour).Unix(), "exp": time.Now().Add(time.Hour * 2).Unix()}),
			wantErr: ErrInvalidToken,
		},
		{
			name:  "valid token",
			token: "Bearer " + makeToken("test-secret", jwt.MapClaims{"iss": "cloud-app", "sub": "user-123", "iat": time.Now().Unix(), "exp": time.Now().Add(time.Hour).Unix()}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := VerifyToken(tt.token)
			if tt.wantErr != nil {
				if err != tt.wantErr {
					t.Errorf("expected %v, got %v", tt.wantErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if claims.Subject != "user-123" {
				t.Errorf("expected subject user-123, got %q", claims.Subject)
			}
		})
	}
}

func makeToken(secret string, claims jwt.MapClaims) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(secret))
	return signed
}
