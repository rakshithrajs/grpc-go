package handlers

import (
	userspb "cloud/gen/users/v1"
	"cloud/internal/handlers"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (u *UsersHandler) GetUserProfile(ctx context.Context, req *userspb.GetUserProfileRequest) (*userspb.GetUserProfileResponse, error) {
	userID, err := handlers.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := u.storage.GetUserByID(ctx, userID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &userspb.GetUserProfileResponse{
		User: &userspb.User{
			Id:    *user.ID,
			Name:  *user.Name,
			Email: *user.Email,
			Phone: *user.Phone,
		},
	}, nil
}
