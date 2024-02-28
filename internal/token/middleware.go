package token

import (
	"errors"
	"strings"
)

func Middle(tokenString string, target any) error {
	if tokenString == "" {
		return errors.New("TODO:")
	}

	parts := strings.Split(tokenString, " ")
	if len(parts) < 2 || parts[0] != "Bearer" {
		return errors.New("TODO:")
	}

	tokenString = parts[1]
	return JWT.Validate(tokenString, target)
}
