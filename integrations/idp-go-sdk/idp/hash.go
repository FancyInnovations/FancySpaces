package idp

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"strings"

	"golang.org/x/crypto/argon2"
)

// PasswordHash generates a secure hash for the given password using the Argon2 algorithm.
func PasswordHash(password string) string {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return ""
	}

	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	return b64Salt + "." + b64Hash
}

func CheckPassword(password, stored string) (bool, error) {
	parts := strings.Split(stored, ".")
	if len(parts) != 2 {
		return false, errors.New("invalid stored password format")
	}

	// Decode salt
	salt, err := base64.RawStdEncoding.DecodeString(parts[0])
	if err != nil {
		return false, err
	}

	// Decode original hash
	expectedHash, err := base64.RawStdEncoding.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	// Recompute hash using SAME parameters
	hash := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)

	// Constant-time comparison
	if subtle.ConstantTimeCompare(hash, expectedHash) == 1 {
		return true, nil
	}

	return false, nil
}
