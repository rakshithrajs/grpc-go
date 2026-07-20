package handlers

import (
	accountpb "github.com/rakshithrajs/cloud/services/account/gen/account/v1"
	"github.com/rakshithrajs/cloud/services/account/internal/storage"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AccountHandler) GetUserProfile(ctx context.Context, req *accountpb.GetUserProfileRequest) (*accountpb.GetUserProfileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := a.storage.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return &accountpb.GetUserProfileResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &accountpb.GetUserProfileResponse{User: &accountpb.User{
		Id:    *user.ID,
		Name:  *user.Name,
		Email: *user.Email,
		Phone: *user.Phone,
	}}, nil
}
