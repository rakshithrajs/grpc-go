package handlers

import (
	"context"
	"errors"

	UMSpb "github.com/rakshithrajs/cloud/UMS/gen/UMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *UMSHandler) GetUserProfile(ctx context.Context, req *UMSpb.GetUserProfileRequest) (*UMSpb.GetUserProfileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	user, err := a.storage.GetUserByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return &UMSpb.GetUserProfileResponse{}, nil
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &UMSpb.GetUserProfileResponse{User: &UMSpb.User{
		Id:    *user.ID,
		Name:  *user.Name,
		Email: *user.Email,
		Phone: *user.Phone,
	}}, nil
}
