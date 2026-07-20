package handlers

import (
	"context"
	"errors"
	"log/slog"
	"strings"

	accountpb "github.com/rakshithrajs/cloud/services/account/gen/account/v1"
	"github.com/rakshithrajs/cloud/services/account/internal/config"
	"github.com/rakshithrajs/cloud/services/account/internal/models"
	"github.com/rakshithrajs/cloud/services/account/internal/storage"
	"github.com/rakshithrajs/cloud/services/account/internal/utils"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrFailedToLoginUser  = errors.New("failed to login user")
	ErrSomethingWentWrong = errors.New("something went wrong")
)

func (a *AccountHandler) LoginUser(ctx context.Context, req *accountpb.LoginUserRequest) (*accountpb.LoginUserResponse, error) {
	payload := models.LoginUserRequest{
		Email:    utils.NormalizeEmail(req.GetEmail()),
		Password: strings.TrimSpace(req.GetPassword()),
	}

	if err := utils.Validate.Struct(payload); err != nil {
		return &accountpb.LoginUserResponse{}, status.Error(codes.InvalidArgument, strings.Join(utils.Errors(err), "; "))
	}

	user, err := a.storage.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		if errors.Is(err, storage.ErrEmailNotFound) {
			return &accountpb.LoginUserResponse{}, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
		}
		slog.Error(logPrefix(fnLoginUser)+"failed to get user by email", slog.Any("error", err))
		return &accountpb.LoginUserResponse{}, status.Error(codes.Internal, ErrFailedToLoginUser.Error())
	}

	if err := bcrypt.CompareHashAndPassword([]byte(*user.Password), []byte(payload.Password)); err != nil {
		return &accountpb.LoginUserResponse{}, status.Error(codes.Unauthenticated, ErrInvalidCredentials.Error())
	}

	cfg, err := config.GetConfig()
	if err != nil {
		slog.Error(logPrefix(fnLoginUser)+"failed to get config", slog.Any("error", err))
		return &accountpb.LoginUserResponse{}, status.Error(codes.Internal, ErrSomethingWentWrong.Error())
	}

	token, err := config.GenerateJWT(*user, cfg.JWTSecret)
	if err != nil {
		slog.Error(logPrefix(fnLoginUser)+"failed to generate JWT", slog.Any("error", err))
		return &accountpb.LoginUserResponse{}, status.Error(codes.Internal, ErrSomethingWentWrong.Error())
	}

	return &accountpb.LoginUserResponse{Token: token}, nil
}
