package token

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
)

var JWT *jwtAuth

type jwtAuth struct {
	secret string
	expire time.Duration
}

func (j *jwtAuth) Generate(data json.RawMessage) (string, error) {
	expiredAt := time.Now().UTC().Add(j.expire)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = expiredAt.Unix()
	json.Unmarshal(data, &claims)
	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (j *jwtAuth) Validate(tokenString string, target any) error {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		expirationTime := time.Unix(int64(claims["exp"].(float64)), 0)
		if time.Now().After(expirationTime) {
			return fmt.Errorf("token has expired")
		}

		raw, err := json.Marshal(claims)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(raw, target); err != nil {
			return err
		}
	}

	return nil
}

func InitJWT(secret string, tokenExpire time.Duration) {
	JWT = &jwtAuth{
		secret: secret,
		expire: tokenExpire,
	}
}
