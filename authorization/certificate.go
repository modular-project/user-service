package authorization

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
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
	// privateBytes, _ := os.LookupEnv("RGE_RSA")

	// publicBytes, _ := os.LookupEnv("RGE_RSA_PUB")
	privateBytes := `-----BEGIN RSA PRIVATE KEY-----
MIICWgIBAAKBgEHTwniN8HPwPqCNWEDMX7qP1K5f4LOhNr01X7J8U5JIMtfJE3uD
RdmndXOIpoTKJVBen0zw1wWvAikGTZd7RV3DgpZ+xr6znB41DeIGhUtfqt2qFhcp
ODAIzWKep+UFoBdJNfq9bqa4kWctURVmVADz8+0Bf6JHV9mBmAC93eDfAgMBAAEC
gYAjfxlDAOdE3awnz5BjgNGuPJknXrRAqRJnfTyZdslp/FzOV7OVyvgDonWHU4zX
1lnAuQWV69lHS1QS4z88DFEe7qbkURJn7tgMUT+FeOYZRHtXvFJSsoODYx01uA2n
aUqLU2pqp9082EyjwwwJfn5AURcWzAtajYHjoOhc9zPSwQJBAINzS2Xp8ECA7kFB
cLCUsun2SUDdMkXu0werp3X9P5brDAx+Mbj9EWQWCd7ApMhXp6Z0PA9VNNtEdtFQ
Dms+PdkCQQCAMtqdJUvY5dlkg94xuS2qb8o4pOS8pWBRTZgePWludhNHOTRbDX8J
XcA7aaqjFlk/4nDqgCXJwDeU3sBXxol3AkBdqosgbMkofXbIgwP0n5C5jCh4kuWe
1WYEQjmKptFoDcbBJC70HUgGJHoWAvmoVGV/A7ZESrfmQmvUDJKpsmlJAkAtSUme
v6EWgsOT1V11dTPjhFAMSHuhBE6NCfsVm54V7lILE/MhwxfASETy9/XWXLu0bJp0
zEYNCgDYbwPFPhYrAkBO97PiuvAcaBVdWJeoEC2OxM9i2c15MegijU9WCK2wzMQP
ey/zSqB+fXLC57wz7Vhb2SaHx96LJBD77K7lyyz5
-----END RSA PRIVATE KEY-----`
	publicBytes := `-----BEGIN PUBLIC KEY-----
MIGeMA0GCSqGSIb3DQEBAQUAA4GMADCBiAKBgEHTwniN8HPwPqCNWEDMX7qP1K5f
4LOhNr01X7J8U5JIMtfJE3uDRdmndXOIpoTKJVBen0zw1wWvAikGTZd7RV3DgpZ+
xr6znB41DeIGhUtfqt2qFhcpODAIzWKep+UFoBdJNfq9bqa4kWctURVmVADz8+0B
f6JHV9mBmAC93eDfAgMBAAE=
-----END PUBLIC KEY-----`

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
