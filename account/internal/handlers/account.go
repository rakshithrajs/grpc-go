package handlers

import (
	"context"
	"errors"

	accountpb "github.com/rakshithrajs/cloud/services/account/gen/account/v1"
	"github.com/rakshithrajs/cloud/services/account/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const nullString = ""

func logPrefix(fn string) string { return "[" + fn + "]: " }

const (
	fnLoginUser         = "LoginUser"
	fnRegisterUser      = "RegisterUser"
	fnUpdateUserProfile = "UpdateUserProfile"
	fnGetUserProfile    = "GetUserProfile"
	fnUserIDFromContext = "UserIDFromContext"
)

var (
	ErrMissingMetadata = errors.New("missing metadata")
	ErrMissingUserID   = errors.New("missing user_id in metadata")
)

func UserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nullString, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 || userIDs[0] == nullString {
		return nullString, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	}

	return userIDs[0], nil
}

type AccountHandler struct {
	accountpb.UnimplementedAccountServer
	storage storage.UserService
}

func NewAccountHandler(store storage.UserService) *AccountHandler {
	return &AccountHandler{storage: store}
}
