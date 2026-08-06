package handlers

import (
	"context"
	"testing"

	TMSpb "github.com/rakshithrajs/cloud/TMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/TMS/internal/config"
	jwtHelper "github.com/rakshithrajs/cloud/TMS/internal/handlers/helpers"
	"google.golang.org/grpc/codes"
)

func TestGenerateToken(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		noConfig     bool
		expectedCode codes.Code
		expectedErr  string
	}{
		{
			name:         "generate token succeeds",
			userID:       "test-user-id",
			expectedCode: codes.OK,
			expectedErr:  config.NullString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := NewTokenHandler()

			resp, err := handler.GenerateToken(context.Background(), &TMSpb.GenerateTokenRequest{UserID: tt.userID})

			jwtHelper.CheckGRPCStatus(t, err, tt.expectedCode, tt.expectedErr)

			if tt.expectedCode == codes.OK {
				if resp.GetAccessToken() == config.NullString {
					t.Errorf("expected non-empty access token")
				}
			}
		})
	}
}
