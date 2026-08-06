package handlers

import (
	"errors"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rakshithrajs/cloud/TMS/internal/config"
)

const (
	// JWTIssuer is the issuer name included in generated JWTs.
	JWTIssuer = "cloud-app"
	// JWTExpiry is the lifetime of a generated JWT.
	JWTExpiry = time.Hour * 24
	// JWTClaimIssuer is the JWT claim key for the issuer.
	JWTClaimIssuer = "iss"
	// JWTClaimSubject is the JWT claim key for the subject.
	JWTClaimSubject = "sub"
	// JWTClaimIssuedAt is the JWT claim key for the issued-at timestamp.
	JWTClaimIssuedAt = "iat"
	// JWTClaimExpiry is the JWT claim key for the expiration timestamp.
	JWTClaimExpiry = "exp"
)

var (
	// ErrMissingBearerToken is returned when the Authorization header is missing or does not start with "Bearer ".
	ErrMissingBearerToken = errors.New("missing bearer token")
	// ErrInvalidToken is returned when the token is malformed, has an invalid signature, or cannot be verified.
	ErrInvalidToken = errors.New("invalid token")
	// ErrTokenExpired is returned when the token has expired.
	ErrTokenExpired = errors.New("token expired")
	// ErrSomethingWentWrong is returned when an unexpected internal error occurs.
	ErrSomethingWentWrong = errors.New("something went wrong")
)

// Claims holds the verified JWT claims returned by VerifyJWT.
type Claims struct {
	Subject  string
	Issuer   string
	IssuedAt int64
}

// GenerateJWT creates and signs a new access token for the given user ID.
func GenerateJWT(userID string) (string, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return "", ErrSomethingWentWrong
	}

	now := time.Now().UTC()
	claims := jwt.MapClaims{
		JWTClaimIssuer:   JWTIssuer,
		JWTClaimSubject:  userID,
		JWTClaimIssuedAt: now.Unix(),
		JWTClaimExpiry:   now.Add(JWTExpiry).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// VerifyJWT validates the provided access token and returns its claims.
func VerifyJWT(tokenString string) (*Claims, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(tokenString, bearerPrefix) {
		return nil, ErrMissingBearerToken
	}

	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, bearerPrefix))
	if tokenString == config.NullString {
		return nil, ErrInvalidToken
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

	if iss, ok := claims[JWTClaimIssuer].(string); !ok || iss != JWTIssuer {
		return nil, ErrInvalidToken
	}

	userID, ok := claims[JWTClaimSubject].(string)
	if !ok || userID == config.NullString {
		return nil, ErrInvalidToken
	}

	if iat, ok := claims[JWTClaimIssuedAt].(float64); !ok || int64(iat) <= 0 || time.Unix(int64(iat), 0).After(time.Now()) {
		return nil, ErrInvalidToken
	}
	return &Claims{
		Issuer:   claims[JWTClaimIssuer].(string),
		Subject:  claims[JWTClaimSubject].(string),
		IssuedAt: int64(claims[JWTClaimIssuedAt].(float64)),
	}, nil
}
