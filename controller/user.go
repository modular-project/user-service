package controller

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"users-service/model"
)

var (
	ErrEmailNotValid      = errors.New("email not valid")
	ErrPasswordNotValid   = errors.New("password not valid")
	ErrEmailAlreadyInUsed = errors.New("email already in use")
	ErrWrongPassword      = errors.New("wrong password")
	ErrUserNotFound       = errors.New("user not found")
	ErrCodeNotFound       = errors.New("code not found")
	ErrExpiredCode        = errors.New("expired code")
	ErrInvalidCode        = errors.New("invalid code")
)

type UserStorager interface {
	Update(*model.User) error         // Update
	Find(id uint) (model.User, error) // don't return password
	Create(*model.User) error
	ChangePassword(UserID uint, password *string) error
	FindByEmail(*string) (model.User, error) //return Password and ID
	Verify(UserID uint) error
	IsEmployee(userID uint) bool
}

type VerificationStorager interface {
	Find(UserID uint) (model.Verification, error)
	Delete(UserID uint) error
	Create(*model.Verification) error
}

type Mailer interface {
	Confirm(dest string, code string) error
}

type UserService struct {
	user UserStorager
	ver  VerificationStorager
	mail Mailer
}

func NewUserService(u UserStorager, ver VerificationStorager, mail Mailer) UserService {
	return UserService{u, ver, mail}
}

func (st UserService) Data(ID uint) (model.User, error) {
	return st.user.Find(ID)
}

func (st UserService) Verify(userID uint, code string) error {
	if code == "" {
		return ErrNullCode
	}
	ver, err := st.ver.Find(userID)
	if err != nil {
		return err
	}
	if ver.ExpiresAt.Before(time.Now()) {
		return ErrExpiredCode
	}
	if !strings.EqualFold(code, ver.Code) {
		return ErrInvalidCode
	}
	err = st.user.Verify(userID)
	if err != nil {
		return fmt.Errorf("st.user.Verify: %w", err)
	}
	return st.ver.Delete(userID)
}

func (st UserService) GenerateCode(userID uint) error {
	user, err := st.user.Find(userID)
	if err != nil {
		return fmt.Errorf("error at User.Find: %w", err)
	}
	code := generateRandomString(CODESIZE)
	err = st.mail.Confirm(user.Email, code)
	if err != nil {
		return fmt.Errorf("error at send email: %w", err)
	}
	m := model.Verification{
		UserID:    userID,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}
	err = st.ver.Create(&m)
	if err != nil {
		return fmt.Errorf("error at create verification in DB: %w", err)
	}
	return nil
}

func (st UserService) ChangePassword(userID uint, password *string) error {
	if password == nil {
		return ErrNullValue
	}
	if !isPasswordValid(*password) {
		return ErrPasswordNotValid
	}
	pwdB, err := hashAndSalt([]byte(*password))
	if err != nil {
		return fmt.Errorf("error at hashAndSalt password: %w", err)
	}
	pwd := string(pwdB)
	return st.user.ChangePassword(userID, &pwd)
}

func (st UserService) UpdateData(user *model.User) error {
	if user.ID == 0 {
		return ErrUserNotFound
	}
	if st.user.IsEmployee(user.ID) {
		return ErrUnauthorizedUser
	}
	return st.user.Update(user)
}
