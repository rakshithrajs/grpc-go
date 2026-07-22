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
	ErrFailedToRegisterUser = errors.New("failed to register user")
)

func (a *UMSHandler) RegisterUser(ctx context.Context, req *UMSpb.RegisterUserRequest) (*UMSpb.RegisterUserResponse, error) {
	payload := models.RegisterUserRequest{
		Name:            strings.TrimSpace(req.GetName()),
		Email:           utils.NormalizeEmail(req.GetEmail()),
		Password:        strings.TrimSpace(req.GetPassword()),
		Phone:           strings.TrimSpace(req.GetPhone()),
		ConfirmPassword: strings.TrimSpace(req.GetConfirmPassword()),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return &UMSpb.RegisterUserResponse{}, status.Error(codes.InvalidArgument, strings.Join(utils.Errors(err), "; "))
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(payload.Password), bcrypt.DefaultCost)
	if err != nil {
		slog.Error(logPrefix(fnRegisterUser)+"failed to hash password", slog.Any("error", err))
		return &UMSpb.RegisterUserResponse{}, status.Error(codes.Internal, ErrFailedToRegisterUser.Error())
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
			return &UMSpb.RegisterUserResponse{}, status.Error(codes.AlreadyExists, err.Error())
		}
		slog.Error(logPrefix(fnRegisterUser)+"failed to create user", slog.Any("error", err))
		return &UMSpb.RegisterUserResponse{}, status.Error(codes.Internal, ErrFailedToRegisterUser.Error())
	}

	return &UMSpb.RegisterUserResponse{User: &UMSpb.User{
		Id:    *newUser.ID,
		Name:  *newUser.Name,
		Email: *newUser.Email,
		Phone: *newUser.Phone,
	}}, nil
}
