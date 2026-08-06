package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/grpcClient"
	"github.com/rakshithrajs/cloud/UMS/internal/middleware"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

const (

	// FnLoginUser is the log prefix for the login handler.
	FnLoginUser = "LoginUserHandler"

	// FnRegisterUser is the log prefix for the register handler.
	FnRegisterUser = "RegisterUserHandler"

	// FnUpdateUserProfile is the log prefix for the update profile handler.
	FnUpdateUserProfile = "UpdateUserProfileHandler"

	// FnGetUserProfile is the log prefix for the get profile handler.
	FnGetUserProfile = "GetUserProfileHandler"
)

// UserHandler exposes HTTP handlers for user account operations.
type UserHandler struct {
	storage   storage.UserService
	tmsClient *grpcClient.TMSClient
}

// NewUserHandler creates a new UserHandler.
func NewUserHandler(store storage.UserService, tmsClient *grpcClient.TMSClient) *UserHandler {
	return &UserHandler{storage: store, tmsClient: tmsClient}
}

// RegisterRoutes wires the user account routes to the provided router group.
func RegisterRoutes(rg *gin.RouterGroup, h *UserHandler, tmsClient *grpcClient.TMSClient) {
	rg.POST("/register", h.RegisterUserHandler)
	rg.POST("/login", h.LoginUserHandler)
	rg.Use(middleware.AuthMiddleware(tmsClient))
	rg.GET("/profile", h.GetUserProfileHandler)
	rg.PATCH("/update", h.UpdateUserHandler)
}
