package encryption

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

func InitEncryptionService() error {
	serviceType := os.Getenv("ENCRYPTION_SERVICE")
	if serviceType == "" {
		serviceType = "vault"
	}

	var svc EncryptionService
	var err error

	switch serviceType {
	case "vault":
		svc, err = createVaultService()
	default:
		return fmt.Errorf("unsupported encryption service: %s", serviceType)
	}

	if err != nil {
		return err
	}

	GlobalEncryptionService = svc
	return nil
}

func createVaultService() (EncryptionService, error) {
	config := api.DefaultConfig()

	// Prioritize VAULT_HOST as requested, fallback to VAULT_ADDR or default
	vaultHost := os.Getenv("VAULT_HOST")
	if vaultHost != "" {
		config.Address = vaultHost
	} else if addr := os.Getenv("VAULT_ADDR"); addr != "" {
		config.Address = addr
	} else {
		config.Address = "http://127.0.0.1:8200"
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create vault client: %w", err)
	}

	// Dev mode token support
	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		client.SetToken(token)
	} else {
		// Default to root for dev environment consistency if not set
		client.SetToken("root")
	}

	// Key name could be configurable, sticking to "backend-key" for now
	return NewVaultEncryptionService(client, "backend-key"), nil
}

func GetDataSourceKey(id int) string {
	return fmt.Sprintf("data_sources:%d:config", id)
}
