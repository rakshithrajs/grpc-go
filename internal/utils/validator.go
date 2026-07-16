package utils

import (
	"errors"
	"net"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var (
	ErrEmailRequired         = errors.New("email is required")
	ErrInvalidEmail          = errors.New("email is invalid")
	ErrPasswordRequired      = errors.New("password is required")
	ErrInvalidPasswordFormat = errors.New("password should be atleast 8 characters long, contain at least one number, one uppercase letter, one lowercase letter, one special character(!, @, #, $, &, _) and no spaces")
	ErrPhoneRequired         = errors.New("phone number is required")
	ErrInvalidPhoneNumber    = errors.New("phone number should contain digits only and should be 10 digits long")
	ErrNameRequired          = errors.New("name is required")
	ErrFileNameRequired      = errors.New("file name is required")
	ErrFileNameTooLong       = errors.New("file name must be at most 150 characters")
	ErrFileNameInvalid       = errors.New("file name contains invalid characters")
	ErrFileNameHasSpaces     = errors.New("file name cannot contain spaces")
)

var (
	Validate = validator.New()
)

func ValidatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	var (
		hasUpper   bool
		hasLower   bool
		hasNumber  bool
		hasSpecial bool
	)

	for _, c := range password {
		switch {
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsNumber(c):
			hasNumber = true
		case strings.ContainsRune("!@#$&_", c):
			hasSpecial = true
		case unicode.IsSpace(c):
			return false
		}
	}
	return hasLower && hasUpper && hasSpecial && hasNumber
}

func ValidateNotEmpty(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

func ValidateEmailDomain(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	mx, err := net.LookupMX(parts[1])
	return err == nil && len(mx) > 0
}

func ValidateFileName(fl validator.FieldLevel) bool {
	name := fl.Field().String()
	if !ValidateNotEmpty(fl) {
		return false
	}
	if len(name) > 150 {
		return false
	}
	if strings.ContainsAny(name, "/\\") {
		return false
	}
	return true
}

func init() {
	Validate.RegisterValidation("isValidPassword", ValidatePassword)
	Validate.RegisterValidation("isValueEmpty", ValidateNotEmpty)
	Validate.RegisterValidation("isValidEmailDomain", ValidateEmailDomain)
	Validate.RegisterValidation("isValidFileName", ValidateFileName)
}

func FirstError(err error) error {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) || len(verrs) == 0 {
		return err
	}

	e := verrs[0]
	switch e.StructField() {
	case "Email":
		if e.Tag() == "required" {
			return ErrEmailRequired
		}
		if e.Tag() == "isValidEmailDomain" {
			return ErrInvalidEmail
		}
		return ErrInvalidEmail
	case "Password":
		if e.Tag() == "required" {
			return ErrPasswordRequired
		}
		if e.Tag() == "isValidPassword" {
			return ErrInvalidPasswordFormat
		}
		if e.Tag() == "min" || e.Tag() == "max" {
			return ErrInvalidPasswordFormat
		}
		return ErrInvalidPasswordFormat
	case "Phone":
		if e.Tag() == "required" || e.Tag() == "isValueEmpty" {
			return ErrPhoneRequired
		}
		return ErrInvalidPhoneNumber
	case "Name":
		if e.Tag() == "required" || e.Tag() == "isValueEmpty" {
			return ErrNameRequired
		}
	case "FileName":
		if e.Tag() == "required" || e.Tag() == "isValueEmpty" {
			return ErrFileNameRequired
		}
		if e.Tag() == "max" {
			return ErrFileNameTooLong
		}
		return ErrFileNameInvalid
	}
	return e
}
