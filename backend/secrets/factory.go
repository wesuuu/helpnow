package secrets

import (
	"fmt"
	"os"

	"github.com/hashicorp/vault/api"
)

func InitSecretStore() error {
	// We default to Vault, can check env if needed like previous factory
	// For now, assume Vault KV

	config := api.DefaultConfig()

	// Prioritize VAULT_HOST as requested
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
		return fmt.Errorf("failed to create vault client: %w", err)
	}

	if token := os.Getenv("VAULT_TOKEN"); token != "" {
		client.SetToken(token)
	} else {
		client.SetToken("root")
	}

	GlobalSecretStore = NewVaultSecretStore(client)
	return nil
}
