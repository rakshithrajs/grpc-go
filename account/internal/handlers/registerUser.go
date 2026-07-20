package handlers

import (
	accountpb "github.com/rakshithrajs/cloud/services/account/gen/account/v1"
	"github.com/rakshithrajs/cloud/services/account/internal/models"
	"github.com/rakshithrajs/cloud/services/account/internal/storage"
	"github.com/rakshithrajs/cloud/services/account/internal/utils"
	"context"
	"errors"
	"log/slog"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrFailedToRegisterUser = errors.New("failed to register user")
)

func (a *AccountHandler) RegisterUser(ctx context.Context, req *accountpb.RegisterUserRequest) (*accountpb.RegisterUserResponse, error) {
	payload := models.RegisterUserRequest{
		Name:            strings.TrimSpace(req.GetName()),
		Email:           utils.NormalizeEmail(req.GetEmail()),
		Password:        strings.TrimSpace(req.GetPassword()),
		Phone:           strings.TrimSpace(req.GetPhone()),
		ConfirmPassword: strings.TrimSpace(req.GetConfirmPassword()),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return &accountpb.RegisterUserResponse{}, status.Error(codes.InvalidArgument, strings.Join(utils.Errors(err), "; "))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error("[RegisterUser]: failed to hash password", slog.Any("error", err))
		return &accountpb.RegisterUserResponse{}, status.Error(codes.Internal, ErrFailedToRegisterUser.Error())
	}

	password := string(hashedPassword)
	newUser, err := a.storage.CreateUser(ctx, &models.User{
		Name:     &payload.Name,
		Email:    &payload.Email,
		Password: &password,
		Phone:    &payload.Phone,
	})
	if err != nil {
		if errors.Is(err, storage.ErrUserEmailAlreadyExists) || errors.Is(err, storage.ErrPhoneNumberAlreadyExists) {
			return &accountpb.RegisterUserResponse{}, status.Error(codes.AlreadyExists, err.Error())
		}
		slog.Error("[RegisterUser]: failed to create user", slog.Any("error", err))
		return &accountpb.RegisterUserResponse{}, status.Error(codes.Internal, ErrFailedToRegisterUser.Error())
	}

	return &accountpb.RegisterUserResponse{User: &accountpb.User{
		Id:    *newUser.ID,
		Name:  *newUser.Name,
		Email: *newUser.Email,
		Phone: *newUser.Phone,
	}}, nil
}
