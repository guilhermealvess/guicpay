package common

import (
	"crypto/sha256"
	"encoding/hex"
)

func ComputeSHA256Hash(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashSum := hasher.Sum(nil)
	hashString := hex.EncodeToString(hashSum)
	return hashString
}
