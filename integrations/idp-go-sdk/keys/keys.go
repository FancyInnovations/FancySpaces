package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log/slog"
	"os"
)

const (
	E2EKeyPath  = "/Users/oliver/Workspace/FancyInnovations/FancyAnalytics/services/idp/cmd/e2e/keys/"
	ProdKeyPath = "/idp_rsa_keys/"
)

// GetOrGenerateRSAKeys tries to load RSA keys from the specified path. If loading fails, it generates new keys and saves them to the path.
func GetOrGenerateRSAKeys(path string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := LoadPrivateKey(path + "private.pem")
	if err != nil {
		slog.Warn("Failed to load private key, generating new keys", "error", err)

		privateKey, publicKey, err := generateNewKeys()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate new keys: %w", err)
		}

		if err := saveKeysToFiles(path, privateKey, publicKey); err != nil {
			return nil, nil, fmt.Errorf("failed to save new keys to files: %w", err)
		}

		return privateKey, publicKey, nil
	}

	publicKey, err := LoadPublicKey(path + "public.pem")
	if err != nil {
		slog.Warn("Failed to load public key, generating new keys", "error", err)
		privateKey, publicKey, err := generateNewKeys()
		if err != nil {
			return nil, nil, fmt.Errorf("failed to generate new keys: %w", err)
		}

		if err := saveKeysToFiles(path, privateKey, publicKey); err != nil {
			return nil, nil, fmt.Errorf("failed to save new keys to files: %w", err)
		}

		return privateKey, publicKey, nil
	}

	return privateKey, publicKey, nil
}

func LoadPrivateKey(path string) (*rsa.PrivateKey, error) {
	keyBytes, err := os.ReadFile(path + "private.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PrivateKey), nil
}

func LoadPublicKey(path string) (*rsa.PublicKey, error) {
	keyBytes, err := os.ReadFile(path + "public.pem")
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(keyBytes)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM block")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return key.(*rsa.PublicKey), nil
}

func saveKeysToFiles(path string, privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) error {
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return err
	}

	privatePem := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	if err := os.WriteFile(path+"private.pem", pem.EncodeToMemory(privatePem), 0600); err != nil {
		return err
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}

	publicPem := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	if err := os.WriteFile(path+"public.pem", pem.EncodeToMemory(publicPem), 0644); err != nil {
		return err
	}

	return nil
}

func generateNewKeys() (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}
