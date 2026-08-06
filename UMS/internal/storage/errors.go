package storage

import "errors"

var (
	// ErrFailedToCreateUser is returned when creating a user fails.
	ErrFailedToCreateUser = errors.New("failed to create user")

	// ErrFailedToGetUserByID is returned when fetching a user by ID fails.
	ErrFailedToGetUserByID = errors.New("failed to get user by ID")

	// ErrFailedToGetUserByEmail is returned when fetching a user by email fails.
	ErrFailedToGetUserByEmail = errors.New("failed to get user by email")

	// ErrFailedToUpdateUser is returned when updating a user fails.
	ErrFailedToUpdateUser = errors.New("failed to update user")

	// ErrFailedToCreateUserFile is returned when creating a user-file mapping fails.
	ErrFailedToCreateUserFile = errors.New("failed to create user file mapping")

	// ErrFailedToDeleteUserFile is returned when deleting a user-file mapping fails.
	ErrFailedToDeleteUserFile = errors.New("failed to delete user file mapping")

	// ErrFailedToUpdateUserFile is returned when updating a user-file mapping fails.
	ErrFailedToUpdateUserFile = errors.New("failed to update user file mapping")

	// ErrFailedToListUserFiles is returned when listing user files fails.
	ErrFailedToListUserFiles = errors.New("failed to list user files")
)
