package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/constant"
)

func getPublicKey(agentName string, keyID string) (*rsa.PublicKey, error) {
	keyReference := constant.AgentKeyName(agentName, keyID)
	publicKeyBytes, err := azure.DownloadBuffer(constant.PublicKeyContainerName, keyReference)
	if err != nil {
		return nil, err
	}

	decodedPublicKeyPem, _ := pem.Decode(publicKeyBytes)

	// Convert the public key to an x509 public key
	parsedPublicKey, err := x509.ParsePKIXPublicKey(decodedPublicKeyPem.Bytes)
	if err != nil {
		log.Debug("Error parsing public key", err)
		return nil, err
	}
	var publicKey *rsa.PublicKey
	publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		log.Debug("Failed coercing public key to rsa.PublicKey format")
		return nil, errors.New("failed to create public key")
	}
	return publicKey, nil
}

func SignMessage(message *constant.Message) error {
	verifierBody := constant.VerifierString(*message)
	verifierHash := sha256.Sum256([]byte(verifierBody))

	keys := config.GetKeys()

	signature, err := rsa.SignPKCS1v15(rand.Reader, keys.PrivateKey, crypto.SHA256, verifierHash[:])
	if err != nil {
		return err
	}

	message.Signature = hex.EncodeToString(signature)

	return nil
}

func VerifyMessage(message constant.Message) error {
	// Compose the verifier of the message
	verifierBody := constant.VerifierString(message)
	verifierHash := sha256.Sum256(verifierBody)

	publicKey, err := getPublicKey(message.Agent, message.KeyID)
	if err != nil {
		log.Debug("Error retrieving public key", err)
		return err
	}

	signature, err := hex.DecodeString(message.Signature)
	if err != nil {
		log.Debug("Error decoding hex string", err)
		return err
	}

	// Verify the message signature with the public key
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, verifierHash[:], signature)
	if err != nil {
		log.Warn(fmt.Sprintf("Key '%s' from agent '%s' is not signed by a trusted source", message.KeyID, message.Agent))
		return err
	}

	return nil
}
