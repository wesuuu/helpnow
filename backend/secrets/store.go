package secrets

import "context"

// SecretStore defines the interface for storing and retrieving secrets.
type SecretStore interface {
	// Write stores data at the specified path.
	// path should be relative to the secrets engine value (e.g., "data_sources/1")
	Write(ctx context.Context, path string, data map[string]interface{}) error

	// Read retrieves data from the specified path.
	Read(ctx context.Context, path string) (map[string]interface{}, error)

	// Delete removes data at the specified path.
	Delete(ctx context.Context, path string) error
}

var GlobalSecretStore SecretStore
