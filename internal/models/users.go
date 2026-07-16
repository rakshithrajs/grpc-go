package models

import "time"

type User struct {
	ID           *string    `json:"id"`
	Name         *string    `json:"name"`
	Email        *string    `json:"email"`
	Password     *string    `json:"password"`
	Phone        *string    `json:"phone"`
	CreatedAtUTC *time.Time `json:"created_at_utc"`
	UpdatedAtUTC *time.Time `json:"updated_at_utc"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name,omitempty" validate:"omitempty,isValueEmpty"`
	Email    *string `json:"email,omitempty" validate:"omitempty,email,isValidEmailDomain,isValueEmpty"`
	Password *string `json:"password,omitempty" validate:"omitempty,min=8,max=64,isValidPassword"`
	Phone    *string `json:"phone,omitempty" validate:"omitempty,isValueEmpty,len=10"`
}
