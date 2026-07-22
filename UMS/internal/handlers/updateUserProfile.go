package handlers

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	UMSpb "github.com/rakshithrajs/cloud/UMS/gen/UMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/models"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
	"github.com/rakshithrajs/cloud/UMS/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoFieldsToUpdate = errors.New("no fields to update")
)

func (a *UMSHandler) UpdateUserProfile(ctx context.Context, req *UMSpb.UpdateUserProfileRequest) (*UMSpb.UpdateUserProfileResponse, error) {
	userID, err := UserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	setField := func(value string) *string {
		if value == nullString {
			return nil
		}
		return &value
	}

	payload := models.UpdateUserRequest{
		Name:  setField(strings.TrimSpace(req.GetName())),
		Email: setField(utils.NormalizeEmail(req.GetEmail())),
		Phone: setField(strings.TrimSpace(req.GetPhone())),
	}

	if req.GetPassword() != nullString {
		password := strings.TrimSpace(req.GetPassword())
		confirmPassword := strings.TrimSpace(req.GetConfirmPassword())
		payload.Password = &password
		payload.ConfirmPassword = &confirmPassword
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return nil, status.Error(codes.InvalidArgument, strings.Join(utils.Errors(err), "; "))
	}

	if payload.Name == nil && payload.Email == nil && payload.Phone == nil && payload.Password == nil {
		return nil, status.Error(codes.InvalidArgument, ErrNoFieldsToUpdate.Error())
	}

	if payload.Password != nil {
		user, err := a.storage.GetUserByID(ctx, userID)
		if err != nil {
			if errors.Is(err, storage.ErrUserNotFound) {
				return &UMSpb.UpdateUserProfileResponse{}, nil
			}
			return nil, status.Error(codes.Internal, err.Error())
		}

		if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(*payload.Password)); err == nil {
			return nil, status.Error(codes.InvalidArgument, storage.ErrPasswordSameAsOldPassword.Error())
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*payload.Password), bcrypt.DefaultCost)
		if err != nil {
			slog.Error(logPrefix(fnUpdateUserProfile)+"failed to hash password", slog.Any("error", err))
			return nil, status.Error(codes.Internal, storage.ErrFailedToUpdateUser.Error())
		}
		hashed := string(hashedPassword)
		payload.Password = &hashed
	}

	if err := a.storage.UpdateUser(ctx, userID, payload); err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return &UMSpb.UpdateUserProfileResponse{}, nil
		}
		if errors.Is(err, storage.ErrUserEmailAlreadyExists) || errors.Is(err, storage.ErrPhoneNumberAlreadyExists) {
			return nil, status.Error(codes.AlreadyExists, err.Error())
		}
		slog.Error(logPrefix(fnUpdateUserProfile)+"failed to update user", slog.Any("error", err))
		return nil, status.Error(codes.Internal, storage.ErrFailedToUpdateUser.Error())
	}

	return &UMSpb.UpdateUserProfileResponse{}, nil
}
