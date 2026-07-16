package handlers

import (
	authpb "cloud/gen/auth/v1"
	"cloud/internal/config"
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

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFailedToLoginUser  = errors.New("failed to login user")
	ErrSomethingWentWrong = errors.New("something went wrong")
)

func (a *AuthHandler) LoginUser(ctx context.Context, req *authpb.LoginUserRequest) (*authpb.LoginUserResponse, error) {
	payload := &models.LoginUserRequest{
		Email:    req.GetEmail(),
		Password: req.GetPassword(),
	}

	payload.Email = strings.TrimSpace(payload.Email)
	payload.Password = strings.TrimSpace(payload.Password)

	if err := utils.Validate.Struct(payload); err != nil {
		if e := utils.FirstError(err); errors.Is(e, utils.ErrEmailRequired) {
			slog.Error("[LoginUser]: email required", slog.Any("error", err))
			return &authpb.LoginUserResponse{}, status.Error(codes.InvalidArgument, utils.FirstError(err).Error())
		}
		if e := utils.FirstError(err); errors.Is(e, utils.ErrPasswordRequired) {
			slog.Error("[LoginUser]: password required", slog.Any("error", err))
			return &authpb.LoginUserResponse{}, status.Error(codes.InvalidArgument, utils.FirstError(err).Error())
		}
		slog.Error("[LoginUser]: Validation Error:", slog.Any("error", err.Error()))
		return &authpb.LoginUserResponse{}, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
	}

	user, err := a.storage.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, storage.ErrEmailDoesntExist) {
			slog.Error("[LoginUser]: Invalid credentials:", slog.Any("error", err.Error()))
			return &authpb.LoginUserResponse{}, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
		}
		return &authpb.LoginUserResponse{}, status.Error(codes.Internal, ErrFailedToLoginUser.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.Password)); err != nil {
		slog.Error("[LoginUser]: Invalid credentials:", slog.Any("error", err.Error()))
		return &authpb.LoginUserResponse{}, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
	}

	token, err := config.GenerateJWT(*user)
	if err != nil {
		slog.Error("[LoginUser]: Failed to generate JWT:", slog.Any("error", err.Error()))
		return &authpb.LoginUserResponse{}, status.Error(codes.Internal, ErrSomethingWentWrong.Error())
	}

	return &authpb.LoginUserResponse{
		Token: token,
	}, nil
}
