package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/willhackett/azure-mft/pkg/config"
)

// EncryptString encrypts using the public key of another agent
func EncryptString(agentName string, keyID string, plaintext string) (string, error) {
	publicKey, err := getPublicKey(agentName, keyID)
	if err != nil {
		log.Debug(fmt.Sprintf("Failed to get public key of agent '%s' with key ID '%s'", agentName, keyID))
		return "", err
	}

	hash := sha256.New()
	bytes := []byte(plaintext)
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, bytes, nil)
	if err != nil {
		log.Debug("Failed to encrypt text", err)
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}

// DecryptString decrypts using the private key of this agent
func DecryptString(ciphertext string) (string, error) {
	keys := config.GetKeys()

	hash := sha256.New()
	bytes, err := hex.DecodeString(ciphertext)
	if err != nil {
		return "", err
	}
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, keys.PrivateKey, bytes, nil)
	if err != nil {
		log.Debug("Failed to decrypt ciphertext", err)
		return "", err
	}
	return string(plaintext), nil
}
