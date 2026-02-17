package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"io"
)

type AES struct {
}

func (e *AES) Encrypt(plaintext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := aesGCM.Seal(nil, nonce, plaintext, nil)

	// Store nonce + ciphertext together
	return append(nonce, ciphertext...), nil
}

func (e *AES) Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := aesGCM.NonceSize()
	nonce := ciphertext[:nonceSize]
	encrypted := ciphertext[nonceSize:]

	plaintext, err := aesGCM.Open(nil, nonce, encrypted, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
