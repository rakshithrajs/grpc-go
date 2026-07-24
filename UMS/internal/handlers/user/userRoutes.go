package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/middleware"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

const (

	// functions name for login user handler
	FnLoginUser = "LoginUserHandler"

	// functions name for register user handler
	FnRegisterUser = "RegisterUserHandler"

	// functions name for update user profile handler
	FnUpdateUserProfile = "UpdateUserProfileHandler"

	// functions name for get user profile handler
	FnGetUserProfile = "GetUserProfileHandler"
)

type UserHandler struct {
	storage storage.UserService
}

func NewUserHandler(store storage.UserService) *UserHandler {
	return &UserHandler{storage: store}
}

func RegisterRoutes(rg *gin.RouterGroup, h *UserHandler) {
	rg.POST("/register", h.RegisterUserHandler)
	rg.POST("/login", h.LoginUserHandler)
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/profile", h.GetUserProfileHandler)
	rg.PATCH("/update", h.UpdateUserHandler)
}
