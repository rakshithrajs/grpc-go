package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	ErrFileNameRequired = errors.New("file name is required")
	ErrFileNameTooLong  = errors.New("file name must be at most 150 characters")
	ErrFileNameInvalid  = errors.New("file name contains invalid characters")
)

const nullString = ""

var Validate = validator.New()

var fileNameRegex = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)

func init() {
	Validate.RegisterValidation("isValueEmpty", validateNotEmpty)
	Validate.RegisterValidation("isValidFileName", validateFileName)
}

func validateNotEmpty(fl validator.FieldLevel) bool {
	return strings.TrimSpace(fl.Field().String()) != nullString
}

func validateFileName(fl validator.FieldLevel) bool {
	name := strings.TrimSpace(fl.Field().String())
	if name == nullString {
		return false
	}
	if len(name) > 150 {
		return false
	}
	return fileNameRegex.MatchString(name)
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
			return ErrFileNameRequired.Error()
		}
		if tag == "max" {
			return ErrFileNameTooLong.Error()
		}
		return ErrFileNameInvalid.Error()
	}

	return fmt.Sprintf("failed %s validation on %s", tag, field)
}
