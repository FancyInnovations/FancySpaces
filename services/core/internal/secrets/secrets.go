package secrets

import (
	"bytes"
	"crypto/rand"
	"time"
)

type Database interface {
	GetSecret(spaceID, key string) (*Secret, error)
	GetSecrets(spaceID string) ([]*Secret, error)
	CreateSecret(secret *Secret) error
	UpdateSecret(secret *Secret) error
	DeleteSecret(spaceID, key string) error
}

type Encrypter interface {
	Encrypt(plaintext, key []byte) ([]byte, error)
	Decrypt(ciphertext, key []byte) ([]byte, error)
}

type Store struct {
	db        Database
	encrypter Encrypter
	masterKey []byte
}

type Configuration struct {
	Database  Database
	Encrypter Encrypter
	MasterKey []byte
}

func New(config Configuration) *Store {
	return &Store{
		db:        config.Database,
		encrypter: config.Encrypter,
		masterKey: config.MasterKey,
	}
}

// GetSecret retrieves the secret without decrypting its value. Use GetDecryptedSecret to get the decrypted value.
func (s *Store) GetSecret(spaceID, key string) (*Secret, error) {
	return s.db.GetSecret(spaceID, key)
}

// GetSecrets retrieves all secrets for a space without decrypting their values. Use GetDecryptedSecret to get the decrypted value of each secret.
func (s *Store) GetSecrets(spaceID string) ([]*Secret, error) {
	return s.db.GetSecrets(spaceID)
}

// GetDecryptedSecret retrieves the secret and decrypts its value before returning it
func (s *Store) GetDecryptedSecret(spaceID, key string) (string, error) {
	secret, err := s.db.GetSecret(spaceID, key)
	if err != nil {
		return "", err
	}

	dek, err := s.encrypter.Decrypt(secret.DataEncryptionKey, s.masterKey)
	if err != nil {
		return "", err
	}

	decryptedValue, err := s.encrypter.Decrypt(secret.Value, dek)
	if err != nil {
		return "", err
	}

	return string(decryptedValue), nil
}

// CreateSecret creates a new secret
func (s *Store) CreateSecret(spaceID string, key, value, description string) error {
	_, err := s.db.GetSecret(spaceID, key)
	if err == nil {
		return ErrSecretAlreadyExists
	}

	dek, err := s.generateDEK()
	if err != nil {
		return err
	}

	encryptedValue, err := s.encrypter.Encrypt([]byte(value), dek)
	if err != nil {
		return err
	}

	encryptedDEK, err := s.encrypter.Encrypt(dek, s.masterKey)
	if err != nil {
		return err
	}

	secret := &Secret{
		SpaceID:           spaceID,
		Key:               key,
		Value:             encryptedValue,
		DataEncryptionKey: encryptedDEK,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		Description:       description,
	}

	return s.db.CreateSecret(secret)
}

// UpdateSecret updates an existing secret. If the value has changed, it generates a new data encryption key and encrypts the new value before updating.
func (s *Store) UpdateSecret(secret *Secret) error {
	existing, err := s.db.GetSecret(secret.SpaceID, secret.Key)
	if err != nil {
		return err
	}

	if bytes.Equal(existing.Value, secret.Value) {
		// Generate a new data encryption key and encrypt the value before updating

		dek, err := s.generateDEK()
		if err != nil {
			return err
		}

		encryptedValue, err := s.encrypter.Encrypt(secret.Value, dek)
		if err != nil {
			return err
		}

		encryptedDEK, err := s.encrypter.Encrypt(dek, s.masterKey)
		if err != nil {
			return err
		}

		secret.Value = encryptedValue
		secret.DataEncryptionKey = encryptedDEK
	}

	secret.UpdatedAt = time.Now()

	return s.db.UpdateSecret(secret)
}

// DeleteSecret deletes a secret by its key
func (s *Store) DeleteSecret(spaceID, key string) error {
	return s.db.DeleteSecret(spaceID, key)
}

// generateDEK generates a random 256-bit data encryption key
func (s *Store) generateDEK() ([]byte, error) {
	key := make([]byte, 32) // 32 bytes = 256 bits
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}
