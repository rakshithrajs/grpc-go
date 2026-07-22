package models

import "time"

type User struct {
	ID           *string    `json:"id"`
	Name         *string    `json:"name"`
	Email        *string    `json:"email"`
	Password     *string    `json:"password"`
	Phone        *string    `json:"phone"`
	CreatedAtUTC *time.Time `json:"createdAtUTC"`
	UpdatedAtUTC *time.Time `json:"updatedAtUTC"`
}

type UpdateUserRequest struct {
	Name            *string `json:"name,omitempty" validate:"omitempty,isValueEmpty,isValidName,max=100"`
	Email           *string `json:"email,omitempty" validate:"omitempty,email,min=5,max=254,isValidEmail"`
	Password        *string `json:"password,omitempty" validate:"omitempty,min=8,max=64,isValidPassword"`
	ConfirmPassword *string `json:"confirm_password,omitempty" validate:"required_with=Password,eqfield=Password"`
	Phone           *string `json:"phone,omitempty" validate:"omitempty,isValueEmpty,isValidPhone"`
}
