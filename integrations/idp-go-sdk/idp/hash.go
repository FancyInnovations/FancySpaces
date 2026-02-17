package idp

import (
	"crypto/sha256"
	"fmt"
)

func PasswordHash(password string) string {
	h := sha256.New()
	h.Write([]byte(password))
	bs := h.Sum(nil)
	return fmt.Sprintf("%x", bs)
}
