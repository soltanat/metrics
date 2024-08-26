package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
	"os"
)

func generateRSAKeys() ([]byte, []byte, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 8192)
	if err != nil {
		return nil, nil, err
	}

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	privateKeyBytes := pem.EncodeToMemory(privateKeyPEM)

	publicKey := &privateKey.PublicKey
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)
	publicKeyPEM := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicKeyBytes = pem.EncodeToMemory(publicKeyPEM)

	return privateKeyBytes, publicKeyBytes, nil
}

func main() {
	privateKeyBytes, publicKeyBytes, err := generateRSAKeys()
	if err != nil {
		log.Fatal(err)
	}

	privateKeyFile, err := os.Create("./private_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer privateKeyFile.Close()
	_, err = privateKeyFile.Write(privateKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	publicKeyFile, err := os.Create("./public_key.pem")
	if err != nil {
		log.Fatal(err)
	}
	defer publicKeyFile.Close()
	_, err = publicKeyFile.Write(publicKeyBytes)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("RSA keys generated successfully.")
}
