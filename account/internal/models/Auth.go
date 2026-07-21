package models

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
