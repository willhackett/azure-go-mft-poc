package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	"github.com/willhackett/azure-mft/pkg/azure"
	"github.com/willhackett/azure-mft/pkg/config"
	"github.com/willhackett/azure-mft/pkg/constant"
)

func generatePrivateKey() *rsa.PrivateKey {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
			fmt.Printf("Cannot generate RSA key\n")
			os.Exit(1)
	}
	return privateKey
}

func dumpPrivateKeyToFile(keysDir string, privateKey *rsa.PrivateKey) error {
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privateKeyBytes,
	}
	privatePem, err := os.Create(keysDir + "/private.pem")
	if err != nil {
			return err
	}
	err = pem.Encode(privatePem, privateKeyBlock)
	if err != nil {
			return err
	}
	return nil
}

func dumpPublicKeyToFile(keysDir string, privateKey *rsa.PrivateKey) error {
	publicKey := &privateKey.PublicKey
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
			return err
	}
	publicKeyBlock := &pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
	}
	publicPem, err := os.Create(keysDir + "/public.pem")
	if err != nil {
			return err
	}
	err = pem.Encode(publicPem, publicKeyBlock)
	if err != nil {
			return err
	}
	return nil
}

func loadKeysFromFile(keysDir string) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKeyPem, err := ioutil.ReadFile(keysDir + "/private.pem")
	if err != nil {
		return nil, nil, err
	}
	publicKeyPem, err := ioutil.ReadFile(keysDir + "/public.pem")
	if err != nil {
		return nil, nil, err
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPem)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPem)
	publicKeyInt, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, nil, err
	}
	publicKey, ok := publicKeyInt.(*rsa.PublicKey)
	if !ok {
		return nil, nil, errors.New("failed to create public key")
	}
	return privateKey, publicKey, nil
}

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, 0755)
		if err != nil {
			return err
		}
	}
	return nil
}

func generatePublicKeyID(publicKey *rsa.PublicKey) (string, error){
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return "", err
	}
	publicKeyHash := sha256.Sum256(publicKeyBytes)
	
	return constant.AgentKeyName(config.GetConfig().Agent.Name, fmt.Sprintf("%x", publicKeyHash)[0:9]), nil
}

func marshalPublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	if err != nil {
			return nil, err
	}

	return pem.EncodeToMemory(publicKeyBlock), nil
}

func getKeys(keysDir string) (*rsa.PrivateKey, *rsa.PublicKey, string, error) {
	err := createDirIfNotExist(keysDir)
	if err != nil {
		return nil, nil, "", err
	}

	privateKey, publicKey, err := loadKeysFromFile(keysDir)
	if err != nil {
		privateKey = generatePrivateKey()
		err = dumpPrivateKeyToFile(keysDir, privateKey)
		if err != nil {
			return nil, nil, "", err
		}
		err = dumpPublicKeyToFile(keysDir, privateKey)
		if err != nil {
			return nil, nil, "", err
		}
	}

	keyID, err := generatePublicKeyID(publicKey)
	if err != nil {
		return nil, nil, "", err
	}

	return privateKey, publicKey, keyID, err
}


func Init() {
	keysDir := config.GetConfig().Paths.KeysDir

	privateKey, publicKey, keyID, err := getKeys(keysDir)
	cobra.CheckErr(err)

	config.SetKeys(privateKey, publicKey, keyID)

	publicKeyBytes, err := marshalPublicKey(publicKey)
	cobra.CheckErr(err)
	err = azure.UploadBuffer(constant.PublicKeyContainerName, keyID, publicKeyBytes)
	cobra.CheckErr(err)
	fmt.Println("Public key updated:", keyID)
}