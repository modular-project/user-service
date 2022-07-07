package controller

import (
	"encoding/base64"
	"math/rand"
	"net/mail"
	"time"
	"unicode"
	"unsafe"
	"users-service/pkg"

	"golang.org/x/crypto/bcrypt"
)

const (
	CHARSET     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	CHARSETBITS = 6                  // 6 bits to represent a letter index
	CHARSETMASK = 1<<CHARSETBITS - 1 // All 1-bits, as many as letterIdxBits
	CHARSETMAX  = 63 / CHARSETBITS   // # of letter indices fitting in 63 bits
	CODESIZE    = 6
)

var (
	ErrRefreshTokenNotFound = pkg.UnauthorizedErr("refresh token not found")
	ErrInvalidRefreshToken  = pkg.UnauthorizedErr("invalid refresh token")
	ErrUnauthorizedUser     = pkg.ForbiddenErr("unauthorized user")
	ErrAlreadyEmployee      = pkg.BadErr("user is already an employee")
	ErrIsNotAnEmployee      = pkg.BadErr("user is not an employee")
	ErrUserIsNotVerified    = pkg.BadErr("user is not verified")
	ErrNullValue            = pkg.BadErr("null value")
	ErrInvalidSalary        = pkg.BadErr("invalid salary")
	ErrEstablishNecesary    = pkg.BadErr("an establishment is necesarry")
	ErrCannotBeAssigned     = pkg.BadErr("cannot be assigned to the establishment")
)

// Generate a random string with size
func generateRandomString(size int) string {
	src := rand.NewSource(time.Now().UnixNano())
	b := make([]byte, size)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := size-1, src.Int63(), CHARSETMAX; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), CHARSETMAX
		}
		if idx := int(cache & CHARSETMASK); idx < len(CHARSET) {
			b[i] = CHARSET[idx]
			i--
		}
		cache >>= CHARSETBITS
		remain--
	}

	return *(*string)(unsafe.Pointer(&b))
}

// hasAndSalt encrypt using RSA
func hashAndSalt(pwd []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

// isEmailValid return true if the email is valid, else return false
func isEmailValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

// isPasswordValid return true if the password is valid
func isPasswordValid(pwd string) bool {
	if len(pwd) < 8 {
		return false
	}

	var (
		hasUpperCase bool
		hasSpecial   bool
		hasNumber    bool
		hasLower     bool
	)

	for _, v := range pwd {
		if hasLower && hasNumber && hasSpecial && hasUpperCase {
			return true
		}
		switch {
		case unicode.IsLower(v):
			hasLower = true
		case unicode.IsUpper(v):
			hasUpperCase = true
		case unicode.IsNumber(v):
			hasNumber = true
		case unicode.IsPunct(v) || unicode.IsSymbol(v):
			hasSpecial = true
		}
	}

	return hasLower && hasNumber && hasSpecial && hasUpperCase
}

func generateFgp(n int) (string, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
