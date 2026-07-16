package handlers

import (
	userspb "cloud/gen/users/v1"
	"cloud/internal/handlers"
	"cloud/internal/models"
	"cloud/internal/storage"
	"cloud/internal/utils"
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

func (u *UsersHandler) UpdateUserProfile(ctx context.Context, req *userspb.UpdateUserProfileRequest) (*userspb.UpdateUserProfileResponse, error) {
	userID, err := handlers.UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	setField := func(value string) *string {
		if value == "" {
			return nil
		}
		return &value
	}

	var payload models.UpdateUserRequest
	payload.Name = setField(req.Name)
	payload.Email = setField(req.Email)
	payload.Phone = setField(req.Phone)
	payload.Password = setField(req.Password)

	if err := utils.Validate.Struct(payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, utils.FirstError(err).Error())
	}

	if payload.Name == nil && payload.Email == nil && payload.Phone == nil && payload.Password == nil {
		return nil, status.Error(codes.InvalidArgument, ErrNoFieldsToUpdate.Error())
	}

	if err = u.storage.UpdateUser(ctx, userID, payload); err != nil {
		if errors.Is(err, storage.ErrUserIDDoesntExist) {
			return &userspb.UpdateUserProfileResponse{}, nil
		}
		if errors.Is(err, storage.ErrPasswordSameAsOldPassword) {
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
		if errors.Is(err, storage.ErrUserEmailAlreadyExists) || errors.Is(err, storage.ErrPhoneNumberAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &userspb.UpdateUserProfileResponse{}, nil
}
