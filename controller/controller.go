package controller

import (
	"net/mail"
	"unicode"
)

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
