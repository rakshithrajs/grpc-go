package models

type RegisterUserRequest struct {
	Name     string `json:"name" validate:"required,isValueEmpty"`
	Email    string `json:"email" validate:"required,email,isValidEmailDomain"`
	Password string `json:"password" validate:"required,min=8,max=64,isValidPassword"`
	Phone    string `json:"phone" validate:"required,isValueEmpty,len=10"`
}

type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email,isValidEmailDomain,isValueEmpty"`
	Password string `json:"password" validate:"required,min=8,max=64,isValidPassword"`
}
