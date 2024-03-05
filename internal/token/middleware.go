package token

import (
	"errors"
	"strings"
)

func Middleware(tokenString string, target any) error {
	if tokenString == "" {
		return errors.New("token invalid")
	}

	parts := strings.Split(tokenString, " ")
	if len(parts) < 2 || parts[0] != "Bearer" {
		return errors.New("token invalid")
	}

	tokenString = parts[1]
	return JWT.Validate(tokenString, target)
}
