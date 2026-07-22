package handlers

import (
	"context"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	UMSpb "github.com/rakshithrajs/cloud/UMS/gen/UMS/v1"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"github.com/rakshithrajs/cloud/UMS/internal/middleware"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func logPrefix(fn string) string { return "[" + fn + "]: " }

const (
	fnLoginUser         = "LoginUser"
	fnRegisterUser      = "RegisterUser"
	fnUpdateUserProfile = "UpdateUserProfile"
	fnGetUserProfile    = "GetUserProfile"
	fnUserIDFromContext = "UserIDFromContext"
)

var (
	ErrMissingMetadata      = errors.New("missing metadata")
	ErrMissingUserID        = errors.New("missing user_id in metadata")
	ErrInvalidJSON          = errors.New("invalid JSON payload")
	ErrFailedToRegisterUser = errors.New("failed to register user")
	ErrInvalidCredentials   = errors.New("invalid credentials")
	ErrFailedToLoginUser    = errors.New("failed to login user")
	ErrSomethingWentWrong   = errors.New("something went wrong")
	ErrNoFieldsToUpdate     = errors.New("no fields to update")
	ErrInvalidID            = errors.New("invalid ID")
	ErrUnauthorized         = errors.New("unauthorized")
)

type UMSHandler struct {
	UMSpb.UnimplementedUMSServer
	storage storage.UserService
}

func NewUMSHandler(store storage.UserService) *UMSHandler {
	return &UMSHandler{storage: store}
}

func RegisterRoutes(rg *gin.RouterGroup, h *UMSHandler) {
	rg.POST("/register", h.RegisterUserHandler)
	rg.POST("/login", h.LoginUserHandler)
	rg.Use(middleware.AuthMiddleware())
	rg.PATCH("/update", h.UpdateUserHandler)
}

func UserIDFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return config.NullString, status.Error(codes.Unauthenticated, ErrMissingMetadata.Error())
	}

	userIDs := md.Get("user_id")
	if len(userIDs) == 0 || userIDs[0] == config.NullString {
		return config.NullString, status.Error(codes.Unauthenticated, ErrMissingUserID.Error())
	}

	return userIDs[0], nil
}

func GetUser(c *gin.Context, name string) (string, error) {
	userID, exists := c.Get("userID")
	if !exists {
		return "", ErrUnauthorized
	}
	return userID.(string), nil
}

func HandleDomainError(c *gin.Context, err error) bool {
	status, ok := domainErrorStatus(err)
	if !ok {
		return false
	}
	c.JSON(status, gin.H{config.ErrorKey: err.Error()})
	return true
}

func domainErrorStatus(err error) (int, bool) {
	switch {
	case errors.Is(err, storage.ErrUserEmailAlreadyExists),
		errors.Is(err, storage.ErrPhoneNumberAlreadyExists):
		return http.StatusConflict, true
	}
	return 0, false
}
