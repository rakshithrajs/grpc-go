package models

import "time"

type User struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Password     string    `json:"password,omitempty"`
	Phone        string    `json:"phone"`
	CreatedAtUTC time.Time `json:"createdAtUTC"`
	UpdatedAtUTC time.Time `json:"updatedAtUTC"`
}

type UpdateUserRequest struct {
	Name            string `json:"name,omitempty" validate:"omitempty,isValueEmpty,isValidName,max=100"`
	Email           string `json:"email,omitempty" validate:"omitempty,email,min=5,max=254,isValidEmail"`
	Password        string `json:"password,omitempty" validate:"omitempty,min=8,max=64,isValidPassword"`
	ConfirmPassword string `json:"confirmPassword,omitempty" validate:"required_with=Password,eqfield=Password"`
	Phone           string `json:"phone,omitempty" validate:"omitempty,isValueEmpty,isValidPhone"`
}

type RegisterUserRequest struct {
	Name            string `json:"name" validate:"required,isValueEmpty,isValidName,max=100"`
	Email           string `json:"email" validate:"required,email,min=5,max=254,isValidEmail"`
	Password        string `json:"password" validate:"required,min=8,max=64,isValidPassword"`
	Phone           string `json:"phone" validate:"required,isValueEmpty,isValidPhone"`
	ConfirmPassword string `json:"confirmPassword" validate:"required,eqfield=Password"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,min=5,max=254,isValidEmail"`
	Password string `json:"password" validate:"required,min=8,max=64,isValidPassword"`
}

type UserFiles struct {
	FileID   string `json:"fileId"`
	UserID   string `json:"userId"`
	FileName string `json:"fileName"`
}

type RenameFileRequest struct {
	NewName string `json:"newName" validate:"required,isValueEmpty,max=255"`
}
