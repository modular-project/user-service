package authorization

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/gbrlsnchs/jwt"
)

var (
	// _signKey   *rsa.PrivateKey
	// _verifyKey *rsa.PublicKey
	_signer jwt.Signer
	once    sync.Once
)

// LoadFiles .
type signerType int

const (
	RSA512 = iota + 1
)

func LoadCertificates(st signerType) error {
	var err error
	once.Do(func() {
		err = loadCertificates(st)
	})

	return err
}

func loadCertificates(st signerType) error {
	privateBytes, _ := os.LookupEnv("RGE_RSA")

	publicBytes, _ := os.LookupEnv("RGE_RSA_PUB")

	return parseRSA(st, []byte(privateBytes), []byte(publicBytes))
}

func parseRSA(st signerType, privateBytes, publicBytes []byte) error {
	var err error
	privateBlock, _ := pem.Decode([]byte(privateBytes))
	signKey, err := x509.ParsePKCS1PrivateKey(privateBlock.Bytes)
	if err != nil {
		return fmt.Errorf("error at ParsePKCS1PrivateKey: %w", err)
	}
	publicBlock, _ := pem.Decode([]byte(publicBytes))
	publicKey, err := x509.ParsePKIXPublicKey(publicBlock.Bytes)
	if err != nil {
		return fmt.Errorf("error at ParsePKCS1Publickey: %w", err)
	}
	verifyKey, _ := publicKey.(*rsa.PublicKey)
	switch st {
	case RSA512:
		_signer = jwt.RS512(signKey, verifyKey)
	default:
		return errors.New("invalid type")
	}
	return nil
}
