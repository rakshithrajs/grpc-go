package handlers

import (
	userspb "cloud/gen/users/v1"

	"cloud/internal/storage"
)

type UsersHandler struct {
	userspb.UnimplementedUsersServer
	storage storage.UserService
}

func NewUsersHandler(userService storage.UserService) *UsersHandler {
	return &UsersHandler{
		storage: userService,
	}
}
