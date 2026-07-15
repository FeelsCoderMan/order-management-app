package validation

import (
	"regexp"
	"errors"
)

var userNameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

func ValidateUsername(username string) error {
	if !userNameRegex.MatchString(username) {
		return errors.New("Username contains invalid characters")
	}

	return nil
}
