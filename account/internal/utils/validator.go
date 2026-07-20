package utils

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"golang.org/x/net/idna"
	"golang.org/x/text/unicode/norm"
)

var (
	ErrEmailRequired           = errors.New("email is required")
	ErrInvalidEmail            = errors.New("email is invalid")
	ErrPasswordRequired        = errors.New("password is required")
	ErrInvalidPassword         = errors.New("password must be 8-64 characters long and contain at least one uppercase letter, one lowercase letter, one number, one special character (!@#$&_) and no spaces")
	ErrPhoneRequired           = errors.New("phone number is required")
	ErrInvalidPhoneNumber      = errors.New("phone number must be exactly 10 digits")
	ErrNameRequired            = errors.New("name is required")
	ErrInvalidName             = errors.New("name can only contain letters")
	ErrNameTooLong             = errors.New("name must be at most 100 characters")
	ErrPasswordMismatch        = errors.New("passwords do not match")
	ErrPasswordConfirmRequired = errors.New("password confirmation is required")
)

var Validate = validator.New()

var (
	nameRegex  = regexp.MustCompile(`^[a-zA-Z]+$`)
	phoneRegex = regexp.MustCompile(`^[0-9]{10}$`)
)

func init() {
	Validate.RegisterValidation("isValueEmpty", validateNotEmpty)
	Validate.RegisterValidation("isValidPassword", validatePassword)
	Validate.RegisterValidation("isValidEmail", validateEmail)
	Validate.RegisterValidation("isValidPhone", validatePhone)
	Validate.RegisterValidation("isValidName", validateName)
}

func validateNotEmpty(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != ""
}

func validatePassword(fl validator.FieldLevel) bool {
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

func NormalizeEmail(raw string) string {
	s := strings.TrimSpace(raw)
	s = norm.NFC.String(s)

	at := strings.LastIndex(s, "@")
	if at < 0 || at == len(s)-1 {
		return s
	}

	local, domain := s[:at], s[at+1:]
	domain = strings.ToLower(domain)
	if puny, err := idna.ToASCII(domain); err == nil {
		domain = puny
	}

	return local + "@" + domain
}

func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	const (
		minLocalLength  = 1
		maxLocalLength  = 64
		minDomainLength = 1
		maxDomainLength = 253
	)

	local := parts[0]
	domain := parts[1]
	if len(local) < minLocalLength || len(local) > maxLocalLength {
		return false
	}
	if len(domain) < minDomainLength || len(domain) > maxDomainLength {
		return false
	}

	mx, err := net.LookupMX(domain)
	return err == nil && len(mx) > 0
}

func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

func Errors(err error) []string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) {
		return []string{err.Error()}
	}

	messages := make([]string, 0, len(verrs))
	for _, e := range verrs {
		messages = append(messages, fmt.Sprintf("%s: %s", e.StructField(), fieldErrorMessage(e.StructField(), e.Tag())))
	}

	return messages
}

func fieldErrorMessage(field, tag string) string {
	switch field {
	case "Name":
		if tag == "required" || tag == "isValueEmpty" {
			return ErrNameRequired.Error()
		}
		if tag == "max" {
			return ErrNameTooLong.Error()
		}
		return ErrInvalidName.Error()
	case "Email":
		if tag == "required" || tag == "isValueEmpty" {
			return ErrEmailRequired.Error()
		}
		return ErrInvalidEmail.Error()
	case "Password":
		if tag == "required" {
			return ErrPasswordRequired.Error()
		}
		return ErrInvalidPassword.Error()
	case "ConfirmPassword":
		if tag == "required" {
			return ErrPasswordConfirmRequired.Error()
		}
		return ErrPasswordMismatch.Error()
	case "Phone":
		if tag == "required" || tag == "isValueEmpty" {
			return ErrPhoneRequired.Error()
		}
		return ErrInvalidPhoneNumber.Error()
	}

	return fmt.Sprintf("failed %s validation on %s", tag, field)
}
