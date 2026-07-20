package config

import (
	"github.com/rakshithrajs/cloud/services/account/internal/models"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrTokenExpired = errors.New("token expired")
)

type Claims struct {
	Subject  string
	Issuer   string
	IssuedAt int64
}

func GenerateJWT(user models.User, secret string) (string, error) {
	claims := jwt.MapClaims{
		"iss": "cloud-app",
		"sub": *user.ID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func ParseJWT(tokenString, secret string) (*Claims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}

	iss, ok := claims["iss"].(string)
	if !ok || iss != "cloud-app" {
		return nil, ErrInvalidToken
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return nil, ErrInvalidToken
	}

	iat, ok := claims["iat"].(float64)
	if !ok || int64(iat) <= 0 || time.Unix(int64(iat), 0).After(time.Now()) {
		return nil, ErrInvalidToken
	}

	return &Claims{
		Subject:  sub,
		Issuer:   iss,
		IssuedAt: int64(iat),
	}, nil
}
