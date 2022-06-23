package authorization

import (
	"errors"
	"fmt"
	"time"

	"github.com/gbrlsnchs/jwt"
)

const (
	iss string = "GoRaSa"
)

var (
	ErrNullToken = errors.New("null token")
)

type Token struct {
	signer jwt.Signer
}

func NewToken() Token {
	return Token{_signer}
}

// GenerateToken .
func (to Token) Create(uid, utp uint) (string, error) {
	claim := jwt.Options{
		ExpirationTime: time.Now().Add(15 * time.Minute),
		Issuer:         iss,
		Public:         map[string]interface{}{"uid": uid, "utp": utp},
	}
	token, err := jwt.Sign(to.signer, &claim)
	if err != nil {
		return "", fmt.Errorf("error at jwt.Sign: %s", err)
	}

	return token, nil
}

func (to Token) CreateRefresh(id, uid uint, fgp *string) (string, error) {
	claim := jwt.Options{
		ExpirationTime: time.Now().Add(168 * time.Hour),
		Issuer:         iss,
		Public: map[string]interface{}{
			"id":  id,
			"uid": uid,
			"fgp": fgp,
		},
	}
	token, err := jwt.Sign(to.signer, &claim)
	if err != nil {
		return "", err
	}

	return token, nil
}

// ValidateToken .
func (to Token) Validate(t *string) (*jwt.JWT, error) {
	if t == nil {
		return nil, ErrNullToken
	}
	jot, err := jwt.FromString(*t)
	if err != nil {
		return &jwt.JWT{}, fmt.Errorf("error at jwt.FromString(%s) : %w", *t, err)
	}
	err = jot.Verify(to.signer)
	if err != nil {
		return &jwt.JWT{}, fmt.Errorf("error at jot.Verify : %w", err)
	}
	err = jot.Validate(jwt.ExpirationTimeValidator(time.Now()), jwt.IssuerValidator(iss), jwt.AlgorithmValidator(jwt.MethodRS512))
	return jot, err
}

// func HashFgp(fgp []byte) []byte {
// 	h := sha256.New()
// 	h.Write(fgp)
// 	return h.Sum(nil)
// }

// // GenerateFgp return a random Fgp string and bytes
// func GenerateFgp(n int) (string, error) {
// 	b := make([]byte, n)
// 	_, err := rand.Read(b)
// 	if err != nil {
// 		return "", err
// 	}
// 	return base64.URLEncoding.EncodeToString(b), nil
// }

// func EqualFpgAndHash(fgp []byte, hash *string) bool {
// 	hashFgp := HashFgp(fgp)
// 	return bytes.Equal(hashFgp, []byte(*hash))
// }
