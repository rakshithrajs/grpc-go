package config

import (
	"errors"
	"time"

	"github.com/rakshithrajs/cloud/UMS/internal/models"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

func GenerateJWT(user models.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"iss": "cloud-app",
		"sub": user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
