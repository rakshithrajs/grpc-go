package models

import "time"

// User represents a registered user account.
type User struct {
	ID           string    `json:"ID"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"-"`
	Phone        string    `json:"phone"`
	CreatedAtUTC time.Time `json:"createdAtUTC"`
	UpdatedAtUTC time.Time `json:"updatedAtUTC"`
}

// UpdateUserRequest contains optional profile fields for updating a user.
type UpdateUserRequest struct {
	Name            string `json:"name,omitempty" validate:"omitempty,isValueEmpty,isValidName,max=100"`
	Email           string `json:"email,omitempty" validate:"omitempty,email,min=5,max=254,isValidEmail"`
	Password        string `json:"password,omitempty" validate:"omitempty,min=8,max=64,isValidPassword"`
	ConfirmPassword string `json:"confirmPassword,omitempty" validate:"required_with=Password,eqfield=Password"`
	Phone           string `json:"phone,omitempty" validate:"omitempty,isValueEmpty,isValidPhone"`
}

// RegisterUserRequest contains the fields required to create a new user account.
type RegisterUserRequest struct {
	Name            string `json:"name" validate:"required,isValueEmpty,isValidName,max=100"`
	Email           string `json:"email" validate:"required,email,min=5,max=254,isValidEmail"`
	Password        string `json:"password" validate:"required,min=8,max=64,isValidPassword"`
	Phone           string `json:"phone" validate:"required,isValueEmpty,isValidPhone"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

// LoginUserRequest contains the credentials used to authenticate a user.
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,min=5,max=254,isValidEmail"`
	Password string `json:"password" validate:"required,min=8,max=64,isValidPassword"`
}
