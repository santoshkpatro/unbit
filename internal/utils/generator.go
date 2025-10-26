package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/oklog/ulid/v2"
)

func GenerateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return base64.StdEncoding.EncodeToString(salt), err
}

func HashPassword(password string, salt string) string {
	hash := sha256.Sum256([]byte(password + salt))
	return base64.StdEncoding.EncodeToString(hash[:])
}

func ComparePassword(password string, salt string, hash string) bool {
	computedHash := HashPassword(password, salt)
	return computedHash == hash
}

func GenerateID(prefix string) string {
	return fmt.Sprintf("%s_%s", prefix, strings.ToLower(ulid.Make().String()))
}
