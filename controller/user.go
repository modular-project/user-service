package controller

import (
	"net/http"
	"strings"
	"time"
	"users-service/model"
	"users-service/pkg"
)

// var (
// 	ErrEmailNotValid      = pkg.AppError{MSG: "email not valid", Code: http.StatusNotFound}
// 	ErrPasswordNotValid   = pkg.AppError{MSG: "password not valid", Code: http.StatusNotFound}
// 	ErrEmailAlreadyInUsed = pkg.AppError{MSG: "email already in use", Code: http.StatusNotFound}
// 	ErrWrongPassword      = pkg.AppError{MSG: "wrong password", Code: http.StatusNotFound}
// 	ErrUserNotFound       = pkg.AppError{MSG: "user not found", Code: http.StatusNotFound}
// 	ErrCodeNotFound       = pkg.AppError{MSG: "code not found", Code: http.StatusNotFound}
// 	ErrExpiredCode        = pkg.AppError{MSG: "expired code", Code: http.StatusNotFound}
// 	ErrInvalidCode        = pkg.AppError{MSG: "invalid code", Code: http.StatusNotFound}
// )

type UserStorager interface {
	Update(*model.User) error         // Update
	Find(id uint) (model.User, error) // don't return password
	Create(*model.User) error
	ChangePassword(UserID uint, password *string) error
	FindByEmail(*string) (model.User, error) //return Password and ID
	Verify(UserID uint) error
	IsEmployee(userID uint) (bool, error)
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
		return pkg.NewAppError("empty code", nil, http.StatusBadRequest)
	}
	ver, err := st.ver.Find(userID)
	if err != nil {
		return pkg.NewAppError("user not found", err, http.StatusBadRequest)
	}
	if ver.ExpiresAt.Before(time.Now()) {
		return pkg.NewAppError("expired code", nil, http.StatusBadRequest)
	}
	if !strings.EqualFold(code, ver.Code) {
		return pkg.NewAppError("invalid code", nil, http.StatusBadRequest)
	}
	err = st.user.Verify(userID)
	if err != nil {
		return pkg.NewAppError("failed to verify user", err, http.StatusInternalServerError)
	}
	if err = st.ver.Delete(userID); err != nil {
		return pkg.NewAppError("failed to remove code", err, http.StatusInternalServerError)
	}
	return nil
}

func (st UserService) GenerateCode(userID uint) error {
	user, err := st.user.Find(userID)
	if err != nil {
		return pkg.NewAppError("user not found", err, http.StatusBadRequest)
	}
	code := generateRandomString(CODESIZE)
	err = st.mail.Confirm(user.Email, code)
	if err != nil {
		return pkg.NewAppError("failed to send email", err, http.StatusInternalServerError)
	}
	m := model.Verification{
		UserID:    userID,
		Code:      code,
		ExpiresAt: time.Now().Add(time.Minute * 15),
	}
	err = st.ver.Create(&m)
	if err != nil {
		return pkg.NewAppError("failed to save the code", err, http.StatusInternalServerError)
	}
	return nil
}

func (st UserService) ChangePassword(userID uint, password *string) error {
	if password == nil {
		return pkg.NewAppError("empty password", nil, http.StatusBadRequest)
	}
	if !isPasswordValid(*password) {
		return pkg.NewAppError("invalid password", nil, http.StatusBadRequest)
	}
	pwdB, err := hashAndSalt([]byte(*password))
	if err != nil {
		return pkg.NewAppError("failed to encrypt password", err, http.StatusInternalServerError)
	}
	pwd := string(pwdB)
	if err = st.user.ChangePassword(userID, &pwd); err != nil {
		return pkg.NewAppError("failed to save password", err, http.StatusInternalServerError)
	}
	return nil
}

func (st UserService) UpdateData(user *model.User) error {
	if user.ID == 0 {
		return pkg.NewAppError("user not found", nil, http.StatusBadRequest)
	}
	is, err := st.user.IsEmployee(user.ID)
	if err != nil {
		return pkg.NewAppError("user not found", err, http.StatusInternalServerError)
	}
	if is {
		return pkg.NewAppError("you have an employee account", nil, http.StatusBadRequest)
	}
	if err := st.user.Update(user); err != nil {
		return pkg.NewAppError("failed to update user", err, http.StatusInternalServerError)
	}
	return nil
}
