package entity

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"
)

type Money int64

const (
	Cent     Money = 1
	Real     Money = 100 * Cent
	MilReais Money = 1000 * Real
)

func (m Money) String() string {
	return fmt.Sprintf("%.2f BRL", float64(m)/100)
}

func (m Money) Absolute() Money {
	if m < 0 {
		return -1 * m
	}

	return m
}

type Password string

func (p *Password) Ok() error {
	return nil
}

func (p *Password) Compare(input string) error {
	parts := strings.Split(string(*p), ":")
	if len(parts) != 3 {
		return errors.New("TODO:")
	}

	method := parts[0]
	salt := parts[1]
	password := parts[2]

	switch method {
	case "SHA256":
		if password != computeSHA256Hash(input+salt) {
			return errors.New("TODO: password invalid")
		}
	}

	return nil
}

func generatePasswordEncoded(password string) Password {
	method := "SHA256"
	salt := computeSHA256Hash(fmt.Sprintf("%d", time.Now().UnixNano()))
	pass := computeSHA256Hash(password + salt)
	return Password(fmt.Sprintf("%s:%s:%s", method, salt, pass))
}

func computeSHA256Hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashSum := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	return hashString
}
