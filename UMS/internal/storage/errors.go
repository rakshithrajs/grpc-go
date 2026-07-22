package storage

import "errors"

var (
	ErrFailedToCreateUser        = errors.New("failed to create user")
	ErrFailedToGetUserByID       = errors.New("failed to get user by ID")
	ErrUserNotFound              = errors.New("user not found")
	ErrUserEmailAlreadyExists    = errors.New("user email already exists")
	ErrPhoneNumberAlreadyExists  = errors.New("phone number already exists")
	ErrFailedToGetUserByEmail    = errors.New("failed to get user by email")
	ErrEmailNotFound             = errors.New("email not found")
	ErrPasswordSameAsOldPassword = errors.New("new password is same as old password")
	ErrFailedToUpdateUser        = errors.New("failed to update user")

	ErrFailedToCreateUserFile = errors.New("failed to create user file mapping")
	ErrFailedToDeleteUserFile = errors.New("failed to delete user file mapping")
	ErrFailedToUpdateUserFile = errors.New("failed to update user file mapping")
	ErrFailedToListUserFiles  = errors.New("failed to list user files")
	ErrUserFileNotFound       = errors.New("user file not found")
	ErrUserFileAlreadyExists  = errors.New("user file mapping already exists")
)
