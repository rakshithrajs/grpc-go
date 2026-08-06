package handlers

import (
	"errors"

	TMSpb "github.com/rakshithrajs/cloud/TMS/gen/TMS/v1"
)

var (
	// ErrFailedToGenerateToken is returned when token generation fails.
	ErrFailedToGenerateToken = errors.New("failed to generate token")
)

// logPrefix returns a formatted string for logging purposes, including the function name.
func logPrefix(fn string) string { return "[" + fn + "]: " }

// TokenHandler implements the TMS gRPC token service.
type TokenHandler struct {
	TMSpb.UnimplementedTokensServer
}

// NewTokenHandler creates a new TokenHandler.
func NewTokenHandler() *TokenHandler {
	return &TokenHandler{}
}
