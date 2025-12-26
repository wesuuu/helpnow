package secrets

import (
	"context"
	"fmt"
	"path"

	"github.com/hashicorp/vault/api"
)

type VaultSecretStore struct {
	client    *api.Client
	mountPath string // defaults to "secret"
}

func NewVaultSecretStore(client *api.Client) *VaultSecretStore {
	return &VaultSecretStore{
		client:    client,
		mountPath: "secret",
	}
}

// Write stores secret data. specific to KV v2, data goes into "data" field.
func (s *VaultSecretStore) Write(ctx context.Context, subPath string, data map[string]interface{}) error {
	// KV v2 write path: secret/data/<subPath>
	fullPath := path.Join(s.mountPath, "data", subPath)

	payload := map[string]interface{}{
		"data": data,
	}

	_, err := s.client.Logical().WriteWithContext(ctx, fullPath, payload)
	if err != nil {
		return fmt.Errorf("failed to write secret at %s: %w", fullPath, err)
	}
	return nil
}

func (s *VaultSecretStore) Read(ctx context.Context, subPath string) (map[string]interface{}, error) {
	// KV v2 read path: secret/data/<subPath>
	fullPath := path.Join(s.mountPath, "data", subPath)

	secret, err := s.client.Logical().ReadWithContext(ctx, fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secret at %s: %w", fullPath, err)
	}

	if secret == nil || secret.Data == nil {
		return nil, nil // Not found
	}

	// KV v2 response structure: secret.Data["data"] holds the actual stored map
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret format at %s", fullPath)
	}

	return data, nil
}

func (s *VaultSecretStore) Delete(ctx context.Context, subPath string) error {
	// KV v2 delete (metadata/destroy): secret/metadata/<subPath> deletes all versions?
	// or specific version.
	// To delete data usually we call metadata endpoint to delete the key entirely.
	fullPath := path.Join(s.mountPath, "metadata", subPath)

	_, err := s.client.Logical().DeleteWithContext(ctx, fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete secret at %s: %w", fullPath, err)
	}
	return nil
}
