package middleware

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
)

var (
	ErrMissingBearerToken = errors.New("missing bearer token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrSomethingWentWrong = errors.New("something went wrong")
)

type Claims struct {
	Subject  string
	Issuer   string
	IssuedAt int64
}

func VerifyToken(tokenString string) (*Claims, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(tokenString, bearerPrefix) {
		return nil, ErrMissingBearerToken
	}

	tokenString = strings.TrimPrefix(tokenString, bearerPrefix)
	if tokenString == config.NullString {
		return nil, ErrMissingBearerToken
	}

	cfg, err := config.GetConfig()
	if err != nil {
		return nil, ErrSomethingWentWrong
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(cfg.JWTSecret), nil
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

	if iss, ok := claims[config.JWTClaimIssuer].(string); !ok || iss != config.JWTIssuer {
		return nil, ErrInvalidToken
	}

	userID, ok := claims[config.JWTClaimSubject].(string)
	if !ok || userID == config.NullString {
		return nil, ErrInvalidToken
	}

	if iat, ok := claims[config.JWTClaimIssuedAt].(float64); !ok || int64(iat) <= 0 || time.Unix(int64(iat), 0).After(time.Now()) {
		return nil, ErrInvalidToken
	}
	return &Claims{
		Issuer:   claims[config.JWTClaimIssuer].(string),
		Subject:  claims[config.JWTClaimSubject].(string),
		IssuedAt: int64(claims[config.JWTClaimIssuedAt].(float64)),
	}, nil
}
