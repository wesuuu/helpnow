package encryption

import "context"

// EncryptionService defines the interface for encrypting and decrypting data.
type EncryptionService interface {
	// Encrypt encrypts the plainText string.
	// distinctId is a unique identifier for the data being encrypted (e.g. table:row:field).
	Encrypt(ctx context.Context, plainText string, distinctId string) (string, error)

	// Decrypt decrypts the cipherText string.
	// distinctId is a unique identifier for the data being decrypted (e.g. table:row:field).
	Decrypt(ctx context.Context, cipherText string, distinctId string) (string, error)
}

var GlobalEncryptionService EncryptionService
