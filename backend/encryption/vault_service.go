package encryption

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/hashicorp/vault/api"
)

type VaultEncryptionService struct {
	client  *api.Client
	keyName string
}

func NewVaultEncryptionService(client *api.Client, keyName string) *VaultEncryptionService {
	return &VaultEncryptionService{
		client:  client,
		keyName: keyName,
	}
}

func (s *VaultEncryptionService) Encrypt(ctx context.Context, plainText string, distinctId string) (string, error) {
	// Vault expects base64 encoded plaintext
	encoded := base64.StdEncoding.EncodeToString([]byte(plainText))

	data := map[string]interface{}{
		"plaintext": encoded,
		"context":   base64.StdEncoding.EncodeToString([]byte(distinctId)),
	}

	path := fmt.Sprintf("transit/encrypt/%s", s.keyName)
	secret, err := s.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return "", fmt.Errorf("failed to encrypt data: %v", err)
	}

	ciphertext, ok := secret.Data["ciphertext"].(string)
	if !ok {
		return "", fmt.Errorf("ciphertext not found in response")
	}

	return ciphertext, nil
}

func (s *VaultEncryptionService) Decrypt(ctx context.Context, cipherText string, distinctId string) (string, error) {
	data := map[string]interface{}{
		"ciphertext": cipherText,
		"context":    base64.StdEncoding.EncodeToString([]byte(distinctId)),
	}

	path := fmt.Sprintf("transit/decrypt/%s", s.keyName)
	secret, err := s.client.Logical().WriteWithContext(ctx, path, data)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt data: %v", err)
	}

	plaintextEncoded, ok := secret.Data["plaintext"].(string)
	if !ok {
		return "", fmt.Errorf("plaintext not found in response")
	}

	plaintext, err := base64.StdEncoding.DecodeString(plaintextEncoded)
	if err != nil {
		return "", fmt.Errorf("failed to decode plaintext: %v", err)
	}

	return string(plaintext), nil
}
