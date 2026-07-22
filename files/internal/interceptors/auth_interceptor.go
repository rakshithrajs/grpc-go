package interceptors

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	ErrMissingAuthHeader  = status.Error(codes.Unauthenticated, "missing authorization header")
	ErrSomethingWentWrong = status.Error(codes.Internal, "something went wrong")
)

func NewAuthInterceptor(UMSnt UMSpb.UMSClient) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, ErrMissingAuthHeader
		}

		authHeaders := md.Get("authorization")
		if len(authHeaders) == 0 {
			return nil, ErrMissingAuthHeader
		}

		UMSCtx := metadata.NewOutgoingContext(ctx, metadata.Pairs("authorization", authHeaders[0]))

		resp, err := UMSClient.GetUserProfile(UMSCtx, &UMSpb.GetUserProfileRequest{})
		if err != nil {
			st, ok := status.FromError(err)
			if ok {
				return nil, st.Err()
			}
			return nil, ErrSomethingWentWrong
		}

		user := resp.GetUser()
		if user == nil || user.GetId() == "" {
			return nil, ErrSomethingWentWrong
		}

		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs("user_id", user.GetId()))

		return handler(ctx, req)
	}
}
