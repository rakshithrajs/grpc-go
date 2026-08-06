package models

import (
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/rakshithrajs/cloud/UMS/internal/config"
	fileModels "github.com/rakshithrajs/cloud/UMS/internal/models"
)

type notEmptyTest struct {
	Value string `validate:"isValueEmpty"`
}

func TestValidateNotEmpty(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid string",
			value:   "hello",
			wantErr: false,
		},
		{
			name:    "valid string with surrounding spaces",
			value:   "  hello  ",
			wantErr: false,
		},
		{
			name:    "empty string",
			value:   config.NullString,
			wantErr: true,
		},
		{
			name:    "only spaces",
			value:   "     ",
			wantErr: true,
		},
		{
			name:    "only tabs and newlines",
			value:   "\t\n",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(notEmptyTest{
				Value: tt.value,
			})

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
		})
	}
}

type passwordTest struct {
	Password string `validate:"isValidPassword"`
}

func TestValidatePassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "valid password",
			password: "Valid123!",
			wantErr:  false,
		},
		{
			name:     "missing uppercase",
			password: "valid123!",
			wantErr:  true,
		},
		{
			name:     "missing lowercase",
			password: "VALID123!",
			wantErr:  true,
		},
		{
			name:     "missing number",
			password: "ValidPass!",
			wantErr:  true,
		},
		{
			name:     "missing special character",
			password: "Valid123",
			wantErr:  true,
		},
		{
			name:     "contains space",
			password: "Valid 123!",
			wantErr:  true,
		},
		{
			name:     "unsupported special character",
			password: "Valid123%",
			wantErr:  true,
		},
		{
			name:     "empty password",
			password: config.NullString,
			wantErr:  true,
		},
		{
			name:     "only spaces",
			password: "     ",
			wantErr:  true,
		},
		{
			name:     "underscore is allowed",
			password: "Valid123_",
			wantErr:  false,
		},
		{
			name:     "ampersand is allowed",
			password: "Valid123&",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(passwordTest{
				Password: tt.password,
			})

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already normalized",
			input:    "user@gmail.com",
			expected: "user@gmail.com",
		},
		{
			name:     "trim surrounding spaces",
			input:    "  user@gmail.com  ",
			expected: "user@gmail.com",
		},
		{
			name:     "convert domain to lowercase",
			input:    "User@GMAIL.COM",
			expected: "User@gmail.com",
		},
		{
			name:     "preserve local part case",
			input:    "JohnDoe@Example.COM",
			expected: "JohnDoe@example.com",
		},
		{
			name:     "unicode domain converted to punycode",
			input:    "user@bücher.de",
			expected: "user@xn--bcher-kva.de",
		},
		{
			name:     "missing at symbol",
			input:    "usergmail.com",
			expected: "usergmail.com",
		},
		{
			name:     "ends with at symbol",
			input:    "user@",
			expected: "user@",
		},
		{
			name:     "multiple at symbols uses last index",
			input:    "a@b@Example.COM",
			expected: "a@b@example.com",
		},
		{
			name:     "empty string",
			input:    config.NullString,
			expected: config.NullString,
		},
		{
			name:     "only spaces",
			input:    "   ",
			expected: config.NullString,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeEmail(tt.input)

			if got != tt.expected {
				t.Fatalf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

type emailTest struct {
	Email string `validate:"isValidEmail"`
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name    string
		email   string
		wantErr bool
	}{
		{
			name:    "valid gmail",
			email:   "ab@gmail.com",
			wantErr: false,
		},
		{
			name:    "valid google",
			email:   "user@google.com",
			wantErr: false,
		},
		{
			name:    "missing at symbol",
			email:   "usergmail.com",
			wantErr: true,
		},
		{
			name:    "multiple at symbols",
			email:   "a@b@gmail.com",
			wantErr: true,
		},
		{
			name:    "local part too short",
			email:   "a@gmail.com",
			wantErr: true,
		},
		{
			name:    "local part too long",
			email:   strings.Repeat("a", 65) + "@gmail.com",
			wantErr: true,
		},
		{
			name:    "domain too short",
			email:   "ab@co",
			wantErr: true,
		},
		{
			name:    "domain too long",
			email:   "ab@" + strings.Repeat("a", 254),
			wantErr: true,
		},
		{
			name:    "domain without mx record",
			email:   "ab@this-domain-should-not-exist-123456789.com",
			wantErr: true,
		},
		{
			name:    "empty string",
			email:   config.NullString,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(emailTest{
				Email: tt.email,
			})

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
		})
	}
}

type phoneTest struct {
	Phone string `validate:"isValidPhone"`
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		name    string
		phone   string
		wantErr bool
	}{
		{
			name:    "valid phone number",
			phone:   "9876543210",
			wantErr: false,
		},
		{
			name:    "less than ten digits",
			phone:   "987654321",
			wantErr: true,
		},
		{
			name:    "more than ten digits",
			phone:   "98765432101",
			wantErr: true,
		},
		{
			name:    "contains letters",
			phone:   "98765abcde",
			wantErr: true,
		},
		{
			name:    "contains special characters",
			phone:   "98765-3210",
			wantErr: true,
		},
		{
			name:    "contains spaces",
			phone:   "98765 3210",
			wantErr: true,
		},
		{
			name:    "empty string",
			phone:   config.NullString,
			wantErr: true,
		},
		{
			name:    "only spaces",
			phone:   "          ",
			wantErr: true,
		},
		{
			name:    "all zeros",
			phone:   "0000000000",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(phoneTest{
				Phone: tt.phone,
			})

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
		})
	}
}

type nameTest struct {
	Name string `validate:"isValidName"`
}

func TestValidateName(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid lowercase",
			value:   "john",
			wantErr: false,
		},
		{
			name:    "valid uppercase",
			value:   "JOHN",
			wantErr: false,
		},
		{
			name:    "valid mixed case",
			value:   "JohnDoe",
			wantErr: false,
		},
		{
			name:    "contains space",
			value:   "John Doe",
			wantErr: true,
		},
		{
			name:    "contains number",
			value:   "John1",
			wantErr: true,
		},
		{
			name:    "contains special character",
			value:   "John!",
			wantErr: true,
		},
		{
			name:    "empty string",
			value:   config.NullString,
			wantErr: true,
		},
		{
			name:    "only spaces",
			value:   "   ",
			wantErr: true,
		},
		{
			name:    "unicode letters",
			value:   "José",
			wantErr: true,
		},
		{
			name:    "hyphenated name",
			value:   "Mary-Jane",
			wantErr: true,
		},
		{
			name:    "apostrophe in name",
			value:   "O'Connor",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(nameTest{
				Name: tt.value,
			})

			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
		})
	}
}

type validationTest struct {
	Email           string `validate:"required,isValueEmpty,min=2,max=254,isValidEmail"`
	Password        string `validate:"required,isValueEmpty,isValidPassword"`
	ConfirmPassword string `validate:"required,eqfield=Password"`
	Phone           string `validate:"required,isValueEmpty,isValidPhone"`
	Name            string `validate:"required,isValueEmpty,max=100,isValidName"`
	NewName         string `validate:"required,isValueEmpty"`
	FileID          string `validate:"required,isValueEmpty,uuid"`
}

func validValidationTest() validationTest {
	return validationTest{
		Email:           "ab@gmail.com",
		Password:        "Valid123!",
		ConfirmPassword: "Valid123!",
		Phone:           "9876543210",
		Name:            "John",
		NewName:         "Johnny",
		FileID:          "550e8400-e29b-41d4-a716-446655440000",
	}
}

type newNameTest struct {
	NewName string `validate:"required,isValueEmpty,max=100"`
}

func TestValidateNewName(t *testing.T) {
	tests := []struct {
		name    string
		newName string
		wantErr bool
	}{
		{
			name:    "new name too long",
			newName: strings.Repeat("a", 101),
			wantErr: true,
		},
		{
			name:    "new name invalid characters",
			newName: "John123",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(newNameTest{NewName: tt.newName})
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
		})
	}
}

func TestFileIDURI(t *testing.T) {
	tests := []struct {
		name    string
		fileID  string
		wantErr bool
	}{
		{
			name:    "valid file id",
			fileID:  "550e8400-e29b-41d4-a716-446655440000",
			wantErr: false,
		},
		{
			name:    "empty file id",
			fileID:  config.NullString,
			wantErr: true,
		},
		{
			name:    "only spaces file id",
			fileID:  "   ",
			wantErr: true,
		},
		{
			name:    "invalid uuid",
			fileID:  "not-a-uuid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(fileModels.FileIDURI{FileID: tt.fileID})
			if (err != nil) != tt.wantErr {
				t.Fatalf("expected error=%v, got %v", tt.wantErr, err != nil)
			}
		})
	}
}

func TestFieldErrors(t *testing.T) {
	tests := []struct {
		name     string
		input    error
		expected map[string]string
	}{
		{
			name: "email required",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = config.NullString
				return v
			}()),
			expected: map[string]string{
				"email": ErrEmailRequired.Error(),
			},
		},
		{
			name: "email too long",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = strings.Repeat("a", 245) + "@gmail.com"
				return v
			}()),
			expected: map[string]string{
				"email": ErrInvalidEmailLength.Error(),
			},
		},
		{
			name: "invalid email",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "ab@invalid-domain-xyz-123.com"
				return v
			}()),
			expected: map[string]string{
				"email": ErrInvalidEmail.Error(),
			},
		},
		{
			name: "password required",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = " "
				v.ConfirmPassword = " "
				return v
			}()),
			expected: map[string]string{
				"password": ErrPasswordRequired.Error(),
			},
		},
		{
			name: "invalid password",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "password"
				v.ConfirmPassword = "password"
				return v
			}()),
			expected: map[string]string{
				"password": ErrInvalidPassword.Error(),
			},
		},
		{
			name: "confirm password required",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = config.NullString
				return v
			}()),
			expected: map[string]string{
				"confirmPassword": ErrPasswordConfirmRequired.Error(),
			},
		},
		{
			name: "password mismatch",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Different123!"
				return v
			}()),
			expected: map[string]string{
				"confirmPassword": ErrPasswordMismatch.Error(),
			},
		},
		{
			name: "phone required",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = config.NullString
				return v
			}()),
			expected: map[string]string{
				"phone": ErrPhoneRequired.Error(),
			},
		},
		{
			name: "invalid phone",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = "123"
				return v
			}()),
			expected: map[string]string{
				"phone": ErrInvalidPhoneNumber.Error(),
			},
		},
		{
			name: "name required",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = "9876543210"
				v.Name = config.NullString
				return v
			}()),
			expected: map[string]string{
				"name": ErrNameRequired.Error(),
			},
		},
		{
			name: "name too long",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = "9876543210"
				v.Name = strings.Repeat("a", 101)
				return v
			}()),
			expected: map[string]string{
				"name": ErrNameTooLong.Error(),
			},
		},
		{
			name: "invalid name",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = "9876543210"
				v.Name = "John1"
				return v
			}()),
			expected: map[string]string{
				"name": ErrInvalidName.Error(),
			},
		},
		{
			name: "new name required",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = "9876543210"
				v.Name = "JohnDoe"
				v.NewName = config.NullString
				return v
			}()),
			expected: map[string]string{
				"newName": ErrNewNameRequired.Error(),
			},
		},
		{
			name: "file id required",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = "9876543210"
				v.Name = "JohnDoe"
				v.NewName = "JaneDoe"
				v.FileID = config.NullString
				return v
			}()),
			expected: map[string]string{
				"fileID": ErrFileIDRequired.Error(),
			},
		},
		{
			name: "file id invalid uuid",
			input: Validate.Struct(func() validationTest {
				v := validValidationTest()
				v.Email = "abc@gmail.com"
				v.Password = "Valid123!"
				v.ConfirmPassword = "Valid123!"
				v.Phone = "9876543210"
				v.Name = "JohnDoe"
				v.NewName = "JaneDoe"
				v.FileID = "invalid-uuid"
				return v
			}()),
			expected: map[string]string{
				"fileID": ErrFileIDInvalidUUID.Error(),
			},
		},
		{
			name: "multiple validation errors",
			input: Validate.Struct(func() validationTest {
				return validationTest{}
			}()),
			expected: map[string]string{
				"email":           ErrEmailRequired.Error(),
				"password":        ErrPasswordRequired.Error(),
				"confirmPassword": ErrPasswordConfirmRequired.Error(),
				"phone":           ErrPhoneRequired.Error(),
				"name":            ErrNameRequired.Error(),
				"newName":         ErrNewNameRequired.Error(),
				"fileID":          ErrFileIDRequired.Error(),
			},
		},
		{
			name:     "non validation error",
			input:    errors.New("some random error"),
			expected: map[string]string{config.NullString: "some random error"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FieldErrors(tt.input)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Fatalf("expected %#v, got %#v", tt.expected, got)
			}
		})
	}
}
