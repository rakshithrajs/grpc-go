package config

import (
	"time"

	"github.com/rakshithrajs/cloud/UMS/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

const (
	JWTIssuer        = "cloud-app"
	JWTExpiry        = time.Hour * 24
	JWTClaimIssuer   = "iss"
	JWTClaimSubject  = "sub"
	JWTClaimIssuedAt = "iat"
	JWTClaimExpiry   = "exp"
)

// GenerateJWT generates a JWT token for the given user using the provided secret key.
func GenerateJWT(user models.User, secret string) (string, error) {
	now := time.Now().UTC()
	claims := jwt.MapClaims{
		JWTClaimIssuer:   JWTIssuer,
		JWTClaimSubject:  user.ID,
		JWTClaimIssuedAt: now.Unix(),
		JWTClaimExpiry:   now.Add(JWTExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
