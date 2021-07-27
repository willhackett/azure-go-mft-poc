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

	message.Signature = hex.EncodeToString(signature)

	return nil
}

func VerifyMessage(message constant.Message) error {
	// Compose the verifier of the message
	verifierBody := constant.VerifierString(message)
	verifierHash := sha256.Sum256(verifierBody)
	
	publicKeyBytes, err := getPublicKey(message.Agent, message.KeyID)
	if err != nil {
		fmt.Println("Error conerting public key")
		return err
	}

	decodedPublicKeyPem, _ := pem.Decode(publicKeyBytes)

	// Convert the public key to an x509 public key
	parsedPublicKey, err := x509.ParsePKIXPublicKey(decodedPublicKeyPem.Bytes)
	if err != nil {
		fmt.Println("Error parsing public key", err)
		return err
	}
	var publicKey *rsa.PublicKey
	publicKey, ok := parsedPublicKey.(*rsa.PublicKey)
	if !ok {
		fmt.Println("Error coerce public key")
		return errors.New("failed to create public key")
	}

	signature, err := hex.DecodeString(message.Signature)
	if err != nil {
		fmt.Println("Error decoding hex string", err)
		return err
	}

	// Verify the message signature with the public key
	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, verifierHash[:], signature)
	if err != nil {
		fmt.Println("Not signed by trusted source")
		return err
	}

	return nil
}