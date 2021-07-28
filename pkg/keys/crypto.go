package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"fmt"

	"github.com/willhackett/azure-mft/pkg/config"
)

// EncryptString encrypts using the public key of another agent
func EncryptString(agentName string, keyID string, plaintext string) (string, error) {
	publicKey, err := getPublicKey(agentName, keyID)
	if err != nil {
		fmt.Println("Failed to get public key", err)
		return "", err
	}

	hash := sha512.New()
	bytes := []byte(plaintext)
	ciphertext, err := rsa.EncryptOAEP(hash, rand.Reader, publicKey, bytes, nil)
	if err != nil {
		return "", err
	}
	return string(ciphertext), nil
}

// DecryptString decrypts using the private key of this agent
func DecryptString(ciphertext string) (string, error) {
	keys := config.GetKeys()

	hash := sha512.New()
	bytes := []byte(ciphertext)
	plaintext, err := rsa.DecryptOAEP(hash, rand.Reader, keys.PrivateKey, bytes, nil)
	if err != nil {
		fmt.Println("Unable to decrypt ciphertext", err)
		return "", err
	}
	return string(plaintext), nil
}
