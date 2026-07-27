package storage

import "errors"

var (
	// message: failed to create user
	ErrFailedToCreateUser = errors.New("failed to create user")

	// message: failed to get user by ID
	ErrFailedToGetUserByID = errors.New("failed to get user by ID")

	// message: failed to get user by email
	ErrFailedToGetUserByEmail = errors.New("failed to get user by email")

	// message: failed to update user
	ErrFailedToUpdateUser = errors.New("failed to update user")

	// message: failed to create user file mapping
	ErrFailedToCreateUserFile = errors.New("failed to create user file mapping")

	// message: failed to delete user file mapping
	ErrFailedToDeleteUserFile = errors.New("failed to delete user file mapping")

	// message: failed to update user file mapping
	ErrFailedToUpdateUserFile = errors.New("failed to update user file mapping")

	// message: failed to list user files
	ErrFailedToListUserFiles = errors.New("failed to list user files")

	// message: user file not found
	ErrUserFileNotFound = errors.New("user file not found")
)
