package grpcClient

import (
	"context"

	TMSpb "github.com/rakshithrajs/cloud/UMS/gen/TMS/v1"
)

// TokenClaims holds the claims extracted from a validated JWT.
type TokenClaims struct {
	UserID   string
	Issuer   string
	IssuedAt int64
}

// GenerateToken requests a new JWT access token for the given user ID from TMS.
func (c *TMSClient) GenerateToken(ctx context.Context, userID string) (string, error) {
	resp, err := c.client.GenerateToken(ctx, &TMSpb.GenerateTokenRequest{UserID: userID})
	if err != nil {
		return "", err
	}
	return resp.GetAccessToken(), nil
}

// ValidateToken verifies an access token with TMS and returns its claims.
func (c *TMSClient) ValidateToken(ctx context.Context, accessToken string) (*TokenClaims, error) {
	resp, err := c.client.ValidateToken(ctx, &TMSpb.ValidateTokenRequest{AccessToken: accessToken})
	if err != nil {
		return nil, err
	}

	claims := resp.GetClaims()
	return &TokenClaims{
		UserID:   claims.GetUserID(),
		Issuer:   claims.GetIssuer(),
		IssuedAt: claims.GetIssuedAt(),
	}, nil
}
