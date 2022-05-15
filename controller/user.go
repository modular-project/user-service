package controller

import (
	"errors"
	"users-service/model"
	"users-service/storage"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailNotValid      = errors.New("email not valid")
	ErrPasswordNotValid   = errors.New("password not valid")
	ErrEmailAlreadyInUsed = errors.New("email already in use")
	ErrWrongPassword      = errors.New("wrong password")
	ErrUserNotFound       = errors.New("user not found")
)

// hasAndSalt encrypt using RSA
func hashAndSalt(pwd []byte) ([]byte, error) {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		return nil, err
	}
	return hash, nil
}

func SignUp(m *model.Account) error {
	if !isEmailValid(m.Email) {
		return ErrEmailNotValid
	}
	if !isPasswordValid(m.Password) {
		return ErrPasswordNotValid
	}
	if storage.DB().Where("email = ?", m.Email).First(&model.Account{}).RowsAffected != 0 {
		return ErrEmailAlreadyInUsed
	}
	pwd, err := hashAndSalt([]byte(m.Password))
	if err != nil {
		return err
	}
	m.Password = string(pwd)
	m.ID = 0
	m.RoleID = 0
	m.IsConfirmated = false
	return storage.DB().Create(m).Error
}

func SignIn(m *model.Account) (model.Account, error) {
	user := model.Account{}
	if !isEmailValid(m.Email) {
		return model.Account{}, ErrEmailNotValid
	}
	if !isPasswordValid(m.Password) {
		return model.Account{}, ErrPasswordNotValid
	}
	rows := storage.DB().First(&user,
		&model.Account{
			Email: m.Email,
		}).RowsAffected
	if rows != 1 {
		return model.Account{}, ErrUserNotFound
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(m.Password)); err != nil {
		return model.Account{}, ErrWrongPassword
	}
	return user, nil
}
