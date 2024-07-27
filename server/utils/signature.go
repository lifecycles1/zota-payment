package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func GenerateSignature(data string) string {
	hasher := sha256.New()
	hasher.Write([]byte(data))
	return hex.EncodeToString(hasher.Sum(nil))
}
