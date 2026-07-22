package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/rakshithrajs/cloud/UMS/internal/middleware"
	"github.com/rakshithrajs/cloud/UMS/internal/storage"
)

const (
	fnLoginUser         = "LoginUserHandler"
	fnRegisterUser      = "RegisterUserHandler"
	fnUpdateUserProfile = "UpdateUserProfileHandler"
	fnGetUserProfile    = "GetUserProfileHandler"
)

type UMSHandler struct {
	storage storage.UserService
}

func NewUMSHandler(store storage.UserService) *UMSHandler {
	return &UMSHandler{storage: store}
}

func RegisterRoutes(rg *gin.RouterGroup, h *UMSHandler) {
	rg.POST("/register", h.RegisterUserHandler)
	rg.POST("/login", h.LoginUserHandler)
	rg.Use(middleware.AuthMiddleware())
	rg.GET("/profile", h.GetUserProfileHandler)
	rg.PATCH("/update/:id", h.UpdateUserHandler)
}
