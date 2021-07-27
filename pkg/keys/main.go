package keys

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
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

func loadKeysFromFile(keysDir string) (error, *rsa.PrivateKey, interface{}) {
	privateKeyPem, err := ioutil.ReadFile(keysDir + "/private.pem")
	if err != nil {
		return err, nil, nil
	}
	publicKeyPem, err := ioutil.ReadFile(keysDir + "/public.pem")
	if err != nil {
		return err, nil, nil
	}

	privateKeyBlock, _ := pem.Decode(privateKeyPem)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return err, nil, nil
	}

	publicKeyBlock, _ := pem.Decode(publicKeyPem)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return err, nil, nil
	}
	return nil, privateKey, publicKey
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

// GetKeys returns keys from the keys directory or generates new ones
func GetKeys(keysDir string) (*rsa.PrivateKey, interface{}, error) {
	err := createDirIfNotExist(keysDir)
	if err != nil {
		return nil, nil, err
	}

	err, privateKey, publicKey := loadKeysFromFile(keysDir)
	if err != nil {
		privateKey = generatePrivateKey()
		err = dumpPrivateKeyToFile(keysDir, privateKey)
		if err != nil {
			return nil, nil, err
		}
		err = dumpPublicKeyToFile(keysDir, privateKey)
		if err != nil {
			return nil, nil, err
		}
	}
	return privateKey, publicKey, err
}
