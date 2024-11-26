package utils

import (
	"errors"
	"regexp"

	"github.com/codepnw/go-auth-cookies/internal/models"
)

type validationResult struct {
	IsValid bool
	Error   error
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func newValidationResult(isValid bool, err error) validationResult {
	return validationResult{
		IsValid: isValid,
		Error:   err,
	}
}

func ValidateEmail(email string) validationResult {
	if email == "" {
		return newValidationResult(false, errors.New("empty email is not allowed"))
	}

	if !emailRegex.MatchString(email) {
		return newValidationResult(false, errors.New("email is invalid"))
	}

	return newValidationResult(true, nil)
}

func ValidatePassword(password string) validationResult {
	if len(password) < 6 {
		return newValidationResult(false, errors.New("password must be of min length 6 chars"))
	}
	return newValidationResult(true, nil)
}

func ValidateUserReq(req models.UserAuthenReq) []string {
	var errs []string

	if v := ValidateEmail(req.Email); !v.IsValid {
		errs = append(errs, v.Error.Error())
	}

	if v := ValidatePassword(req.Password); !v.IsValid {
		errs = append(errs, v.Error.Error())
	}

	return errs
}
 