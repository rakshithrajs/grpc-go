package models

import (
	"errors"
	"net"
	"regexp"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
	"github.com/rakshithrajs/cloud/UMS/internal/config"
	"golang.org/x/net/idna"
	"golang.org/x/text/unicode/norm"
)

var (
	// message: email is required
	ErrEmailRequired = errors.New("email is required")

	// message: email is invalid
	ErrInvalidEmail = errors.New("email is invalid")

	// message: email must be at most 254 characters
	ErrInvalidEmailLength = errors.New("email must be at most 254 characters")

	// message: password is required
	ErrPasswordRequired = errors.New("password is required")

	// message: password must be 8-64 characters long and contain at least one uppercase letter, one lowercase letter, one number, one special character (!@#$&_) and no spaces
	ErrInvalidPassword = errors.New("password must be 8-64 characters long and contain at least one uppercase letter, one lowercase letter, one number, one special character (!@#$&_) and no spaces")

	// message: phone number is required
	ErrPhoneRequired = errors.New("phone number is required")

	// message: phone number must be exactly 10 digits
	ErrInvalidPhoneNumber = errors.New("phone number must be exactly 10 digits")

	// message: name is required
	ErrNameRequired = errors.New("name is required")

	// message: name can only contain letters
	ErrInvalidName = errors.New("name can only contain letters")

	// message: name must be at most 100 characters
	ErrNameTooLong = errors.New("name must be at most 100 characters")

	// message: passwords do not match
	ErrPasswordMismatch = errors.New("passwords do not match")

	// message: password confirmation is required
	ErrPasswordConfirmRequired = errors.New("password confirmation is required")

	// message: new name is required
	ErrNewNameRequired = errors.New("new name is required")

	// message: file ID is required
	ErrFileIDRequired = errors.New("file ID is required")

	// message: file ID has invalid UUID
	ErrFileIDInvalidUUID = errors.New("file ID has invalid UUID")
)

var Validate = validator.New()

// constants for field names used in validation
const (
	fieldEmail           = "Email"
	fieldPassword        = "Password"
	fieldConfirmPassword = "ConfirmPassword"
	fieldPhone           = "Phone"
	fieldName            = "Name"
	fieldNewName         = "NewName"
	fieldFileID          = "FileID"
)

// FileIDPayload is a single-field validation struct for fileID path parameters.
type FileIDPayload struct {
	FileID string `validate:"required,isValueEmpty,uuid" json:"fileID"`
}

// Regular expressions for validating names and phone numbers
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

// validateNotEmpty checks if the field value is not empty after trimming spaces.
func validateNotEmpty(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != config.NullString
}

// validatePassword checks if the password meets the required criteria: at least one uppercase letter, one lowercase letter, one number, one special character (!@#$&_), and no spaces.
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

// NormalizeEmail normalizes the email address by trimming spaces, converting the domain to lowercase, and converting the domain to ASCII.
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

// validateEmail checks if the email address is valid by ensuring it has a proper format and that the domain has MX records.
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()

	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}

	const (
		minLocalLength  = 2
		maxLocalLength  = 64
		minDomainLength = 3
		maxDomainLength = 253
	)

	if len(parts[0]) < minLocalLength || len(parts[0]) > maxLocalLength {
		return false
	}

	domain := parts[1]
	if len(domain) < minDomainLength || len(domain) > maxDomainLength {
		return false
	}

	mx, err := net.LookupMX(domain)
	return err == nil && len(mx) > 0
}

// validatePhone checks if the phone number is exactly 10 digits long.
func validatePhone(fl validator.FieldLevel) bool {
	return phoneRegex.MatchString(fl.Field().String())
}

// validateName checks if the name contains only letters (a-z, A-Z).
func validateName(fl validator.FieldLevel) bool {
	return nameRegex.MatchString(fl.Field().String())
}

// fieldNames maps struct field names to their corresponding JSON field names for error reporting.
var fieldNames = map[string]string{
	fieldEmail:           "email",
	fieldPassword:        "password",
	fieldConfirmPassword: "confirmPassword",
	fieldPhone:           "phone",
	fieldName:            "name",
	fieldNewName:         "newName",
	fieldFileID:          "fileID",
}

// fieldError maps a validator.FieldError to a corresponding error message based on the field and validation tag.
func fieldError(e validator.FieldError) error {
	if e.StructField() == fieldEmail {
		switch e.Tag() {
		case "required", "isValueEmpty":
			return ErrEmailRequired
		case "min", "max":
			return ErrInvalidEmailLength
		default:
			return ErrInvalidEmail
		}
	}
	if e.StructField() == fieldPassword {
		switch e.Tag() {
		case "required", "isValueEmpty":
			return ErrPasswordRequired
		default:
			return ErrInvalidPassword
		}
	}
	if e.StructField() == fieldConfirmPassword {
		switch e.Tag() {
		case "required":
			return ErrPasswordConfirmRequired
		default:
			return ErrPasswordMismatch
		}
	}
	if e.StructField() == fieldPhone {
		switch e.Tag() {
		case "required", "isValueEmpty":
			return ErrPhoneRequired
		default:
			return ErrInvalidPhoneNumber
		}
	}
	if e.StructField() == fieldName {
		switch e.Tag() {
		case "required", "isValueEmpty":
			return ErrNameRequired
		case "max":
			return ErrNameTooLong
		default:
			return ErrInvalidName
		}
	}
	if e.StructField() == fieldNewName {
		switch e.Tag() {
		case "required", "isValueEmpty":
			return ErrNewNameRequired
		case "max":
			return ErrNameTooLong
		default:
			return e
		}
	}
	if e.StructField() == fieldFileID {
		switch e.Tag() {
		case "required", "isValueEmpty":
			return ErrFileIDRequired
		case "uuid":
			return ErrFileIDInvalidUUID
		default:
			return e
		}
	}
	return e
}

// FieldErrors processes the validation errors and returns a map of field names to their corresponding error messages.
func FieldErrors(err error) map[string]string {
	var verrs validator.ValidationErrors
	if !errors.As(err, &verrs) || len(verrs) == 0 {
		return map[string]string{config.NullString: err.Error()}
	}

	seen := map[string]bool{}
	result := map[string]string{}
	for _, e := range verrs {
		name := fieldNames[e.StructField()]
		if name == config.NullString {
			name = strings.ToLower(e.StructField())
		}
		if seen[name] {
			continue
		}
		seen[name] = true
		result[name] = fieldError(e).Error()
	}
	return result
}
