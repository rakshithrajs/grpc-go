package handlers

import (
	authpb "cloud/gen/auth/v1"
	"cloud/internal/storage"
)

type AuthHandler struct {
	authpb.UnimplementedAuthServer
	storage storage.UserService
}

func NewAuthHandler(userService storage.UserService) *AuthHandler {
	return &AuthHandler{
		storage: userService,
	}
}
