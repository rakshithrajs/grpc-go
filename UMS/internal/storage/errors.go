package storage

import "errors"

var (
	ErrFailedToCreateUser     = errors.New("failed to create user")
	ErrFailedToGetUserByID    = errors.New("failed to get user by ID")
	ErrFailedToGetUserByEmail = errors.New("failed to get user by email")
	ErrFailedToUpdateUser     = errors.New("failed to update user")

	ErrFailedToCreateUserFile = errors.New("failed to create user file mapping")
	ErrFailedToDeleteUserFile = errors.New("failed to delete user file mapping")
	ErrFailedToUpdateUserFile = errors.New("failed to update user file mapping")
	ErrFailedToListUserFiles  = errors.New("failed to list user files")
	ErrUserFileNotFound       = errors.New("user file not found")
)
