package handlers

import (
	"context"

	TMSpb "github.com/rakshithrajs/cloud/TMS/gen/TMS/v1"
	jwtHelper "github.com/rakshithrajs/cloud/TMS/internal/handlers/helpers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidateToken verifies an access token and returns its claims.
func (t *TokenHandler) ValidateToken(ctx context.Context, req *TMSpb.ValidateTokenRequest) (*TMSpb.ValidateTokenResponse, error) {
	claims, err := jwtHelper.VerifyJWT(req.GetAccessToken())
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &TMSpb.ValidateTokenResponse{
		Claims: &TMSpb.TokenClaims{
			UserID:   claims.Subject,
			Issuer:   claims.Issuer,
			IssuedAt: claims.IssuedAt,
		},
	}, nil
}
