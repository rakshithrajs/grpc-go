package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingMetadata    = status.Error(codes.Unauthenticated, "missing metadata")
	ErrMissingUserID      = status.Error(codes.Unauthenticated, "missing user id in metadata")
)

func AuthInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrMissingMetadata
		}

		userIDs := md.Get("x-user-id")
		if len(userIDs) == 0 || userIDs[0] == "" {
			return nil, ErrMissingUserID
		}

		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("x-user-id", userIDs[0]))

		return handler(ctx, req)
	}
}
