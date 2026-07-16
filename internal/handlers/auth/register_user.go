package handlers

import (
	authpb "cloud/gen/auth/v1"
	userspb "cloud/gen/users/v1"
	"cloud/internal/models"
	"cloud/internal/storage"
	"cloud/internal/utils"
	"context"
	"errors"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *AuthHandler) RegisterUser(ctx context.Context, req *authpb.RegisterUserRequest) (*authpb.RegisterUserResponse, error) {
	payload := models.RegisterUserRequest{
		Name:     req.GetName(),
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
		Phone:    req.GetPhone(),
	}

	if err := utils.Validate.Struct(&payload); err != nil {
		return &authpb.RegisterUserResponse{}, status.Error(codes.InvalidArgument, utils.FirstError(err).Error())
	}

	cleanedName := strings.TrimSpace(payload.Name)
	cleanedEmail := strings.TrimSpace(payload.Email)
	cleanedPassword := strings.TrimSpace(payload.Password)
	cleanedPhone := strings.TrimSpace(payload.Phone)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cleanedPassword), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("[Register]: failed to hash password", slog.Any("error", err))
		return &authpb.RegisterUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	password := string(hashedPassword)
	newUser, err := a.storage.CreateUser(ctx, &models.User{
		Name:     &cleanedName,
		Email:    &cleanedEmail,
		Password: &password,
		Phone:    &cleanedPhone,
	})
	if err != nil {
		if errors.Is(err, storage.ErrUserEmailAlreadyExists) || errors.Is(err, storage.ErrPhoneNumberAlreadyExists) {
			return &authpb.RegisterUserResponse{}, status.Error(codes.AlreadyExists, err.Error())
		}
		return &authpb.RegisterUserResponse{}, status.Error(codes.Internal, err.Error())
	}

	return &authpb.RegisterUserResponse{User: &userspb.User{
		Id:    *newUser.ID,
		Name:  *newUser.Name,
		Email: *newUser.Email,
		Phone: *newUser.Phone,
	}}, nil
}
