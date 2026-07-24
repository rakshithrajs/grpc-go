package config

import (
	"errors"
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

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

func GenerateJWT(user models.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		JWTClaimIssuer:   JWTIssuer,
		JWTClaimSubject:  user.ID,
		JWTClaimIssuedAt: time.Now().Unix(),
		JWTClaimExpiry:   time.Now().Add(JWTExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
