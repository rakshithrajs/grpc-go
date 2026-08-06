package handlers

import (
	"context"
	"log/slog"

	TMSpb "github.com/rakshithrajs/cloud/TMS/gen/TMS/v1"
	"github.com/rakshithrajs/cloud/TMS/internal/config"
	jwtHelper "github.com/rakshithrajs/cloud/TMS/internal/handlers/helpers"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const fnGenerateToken = "GenerateToken"

// GenerateToken issues a new access token for the requested user.
func (t *TokenHandler) GenerateToken(ctx context.Context, req *TMSpb.GenerateTokenRequest) (*TMSpb.GenerateTokenResponse, error) {
	tokenString, err := jwtHelper.GenerateJWT(req.GetUserID())
	if err != nil {
		slog.Error(logPrefix(fnGenerateToken)+"failed to generate token", slog.Any(config.ErrorKey, err))
		return nil, status.Error(codes.Internal, ErrFailedToGenerateToken.Error())
	}

	return &TMSpb.GenerateTokenResponse{
		AccessToken: tokenString,
	}, nil
}
