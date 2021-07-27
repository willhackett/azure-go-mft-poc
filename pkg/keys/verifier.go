package keys

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"errors"

	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/constant"
)

func getPublicKey(agentName string, keyID string) ([]byte, error) {
	keyReference := constant.AgentKeyName(agentName, keyID)

	return azure.DownloadBuffer(constant.PublicKeyContainerName, keyReference)
}

func SignMessage(message *constant.Message) error {
	verifierBody := constant.VerifierString(*message)
	verifierHash := sha256.Sum256([]byte(verifierBody))

	keys := config.GetKeys()

	signature, err := rsa.SignPKCS1v15(rand.Reader, keys.PrivateKey, crypto.SHA256, verifierHash[:])
	if err != nil {
		return err
	}

	message.Signature = base64.StdEncoding.EncodeToString(signature)

	return nil
}

func VerifyMessage(message constant.Message) error {
	// Compose the verifier of the message
	verifierBody := constant.VerifierString(message)
	verifierHash := sha256.Sum256(verifierBody)
	
	publicKeyBytes, err := getPublicKey(message.Agent, message.KeyID)
	if err != nil {
		return err
	}

	// Convert the public key to an x509 public key
	var publicKey *rsa.PublicKey
	parsedPublicKey, err := x509.ParsePKIXPublicKey(publicKeyBytes)
	if err != nil {
		return err
	}
	publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		return errors.New("failed to create public key")
	}

	signature, err := base64.StdEncoding.DecodeString(message.Signature)
	if err != nil {
		return err
	}

	// Verify the message signature with the public key
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, verifierHash[:], signature)
	if err != nil {
		return err
	}

	return nil
}