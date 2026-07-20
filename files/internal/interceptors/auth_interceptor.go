package interceptors

import (
	"context"
	"errors"
	"log/slog"
	"strings"
	"time"

	accountpb "github.com/rakshithrajs/cloud/services/account/gen/account/v1"
	"github.com/rakshithrajs/cloud/services/files/internal/config"

	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	functionName = "AuthInterceptor"
	logPrefix    = "[" + functionName + "]: "
	nullString   = ""
)

var (
	ErrMissingAuthHeader  = errors.New("missing authorization header")
	ErrMissingBearerToken = errors.New("missing bearer token")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrSomethingWentWrong = errors.New("something went wrong")
)

var publicMethods = map[string]bool{
	accountpb.Account_RegisterUser_FullMethodName: true,
	accountpb.Account_LoginUser_FullMethodName:    true,
}

func AuthInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	if ok := publicMethods[info.FullMethod]; ok {
		return handler(ctx, req)
	}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, ErrMissingAuthHeader.Error())
	}

	authHeaders := md.Get("authorization")
	if len(authHeaders) == 0 {
		slog.Error(logPrefix, slog.String("error", ErrMissingAuthHeader.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrMissingAuthHeader.Error())
	}

	authHeader := authHeaders[0]
	if authHeader == nullString {
		slog.Error(logPrefix, slog.String("error", ErrMissingAuthHeader.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrMissingAuthHeader.Error())
	}

	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		slog.Error(logPrefix, slog.String("error", ErrMissingBearerToken.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrMissingBearerToken.Error())
	}

	tokenString := strings.TrimPrefix(authHeader, bearerPrefix)
	if tokenString == nullString {
		slog.Error(logPrefix, slog.String("error", ErrMissingBearerToken.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrMissingBearerToken.Error())
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(logPrefix, slog.String("error", err.Error()))
		return nil, status.Error(codes.Internal, ErrSomethingWentWrong.Error())
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, ErrInvalidToken
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil || !token.Valid {
		if errors.Is(err, jwt.ErrTokenExpired) {
			slog.Error(logPrefix, slog.String("error", ErrTokenExpired.Error()))
			return nil, status.Error(codes.Unauthenticated, ErrTokenExpired.Error())
		}
		slog.Error(logPrefix, slog.String("error", ErrInvalidToken.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrInvalidToken.Error())
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		slog.Error(logPrefix, slog.String("error", ErrInvalidToken.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrInvalidToken.Error())
	}

	if iss, ok := claims["iss"].(string); !ok || iss != "cloud-app" {
		slog.Error(logPrefix, slog.String("error", ErrInvalidToken.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrInvalidToken.Error())
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == nullString {
		slog.Error(logPrefix, slog.String("error", ErrInvalidToken.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrInvalidToken.Error())
	}

	if iat, ok := claims["iat"].(float64); !ok || int64(iat) <= 0 || time.Unix(int64(iat), 0).After(time.Now()) {
		slog.Error(logPrefix, slog.String("error", ErrInvalidToken.Error()))
		return nil, status.Error(codes.Unauthenticated, ErrInvalidToken.Error())
	}

	ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", userID))

	return handler(ctx, req)
}
