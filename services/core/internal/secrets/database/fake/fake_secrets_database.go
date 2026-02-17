package fake

import (
	"sync"

	"github.com/fancyinnovations/fancyspaces/core/internal/secrets"
)

type DB struct {
	Items map[string]map[string]*secrets.Secret
	mu    sync.RWMutex
}

func NewDB() *DB {
	return &DB{
		Items: make(map[string]map[string]*secrets.Secret),
	}
}

func (db *DB) GetSecret(spaceID, key string) (*secrets.Secret, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	secretsMap, exists := db.Items[spaceID]
	if !exists {
		return nil, secrets.ErrSecretNotFound
	}

	secret, exists := secretsMap[key]
	if !exists {
		return nil, secrets.ErrSecretNotFound
	}

	return secret, nil
}

func (db *DB) GetSecrets(spaceID string) ([]*secrets.Secret, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	secretsMap, exists := db.Items[spaceID]
	if !exists {
		return nil, nil
	}

	var secretsList []*secrets.Secret
	for _, secret := range secretsMap {
		secretsList = append(secretsList, secret)
	}

	return secretsList, nil
}

func (db *DB) CreateSecret(secret *secrets.Secret) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	if _, exists := db.Items[secret.SpaceID]; !exists {
		db.Items[secret.SpaceID] = make(map[string]*secrets.Secret)
	}

	if _, exists := db.Items[secret.SpaceID][secret.Key]; exists {
		return secrets.ErrSecretAlreadyExists
	}

	db.Items[secret.SpaceID][secret.Key] = secret
	return nil
}

func (db *DB) UpdateSecret(secret *secrets.Secret) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	secretsMap, exists := db.Items[secret.SpaceID]
	if !exists {
		return secrets.ErrSecretNotFound
	}

	if _, exists := secretsMap[secret.Key]; !exists {
		return secrets.ErrSecretNotFound
	}

	secretsMap[secret.Key] = secret
	return nil
}

func (db *DB) DeleteSecret(spaceID, key string) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	secretsMap, exists := db.Items[spaceID]
	if !exists {
		return secrets.ErrSecretNotFound
	}

	if _, exists := secretsMap[key]; !exists {
		return secrets.ErrSecretNotFound
	}

	delete(secretsMap, key)
	return nil
}
