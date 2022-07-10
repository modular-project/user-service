package controller

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"
	"users-service/model"
	"users-service/pkg"

	"github.com/gbrlsnchs/jwt"
	"golang.org/x/crypto/bcrypt"
)

type SignStorager interface {
	Create(*model.LogIn) error
	Find(string) (model.LogIn, error)
}

type RefreshStorager interface {
	Delete(id uint) error
	Find(id uint) (model.Refresh, error)
	Create(*model.Refresh) error
}

type Validater interface {
	Validate(*model.LogIn) error
	UType() pkg.UserType
}

type Tokener interface {
	Create(userID, userType uint) (string, error)
	CreateRefresh(userID, userType uint, fgp *string) (string, error)
	Validate(*string) (*jwt.JWT, error)
}

type SignService struct {
	re RefreshStorager
	si SignStorager
	va Validater
	to Tokener
}

func NewSignService(re RefreshStorager, va Validater, si SignStorager, to Tokener) SignService {
	return SignService{re: re, va: va, si: si, to: to}
}

func (ss SignService) SignUp(l *model.LogIn) error {
	err := ss.va.Validate(l)
	if err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	_, err = ss.si.Find(l.User)
	if err != nil && !errors.Is(err, pkg.ErrNoRowsAffected) {
		return pkg.NewAppError("user not found", err, http.StatusInternalServerError)
	}
	if err == nil {
		return pkg.NewAppError("email already in used", nil, http.StatusBadRequest)
	}
	pwd, err := hashAndSalt([]byte(l.Password))
	if err != nil {
		return pkg.NewAppError("fail at encrypt password", err, http.StatusInternalServerError)
	}
	l.Password = string(pwd)

	if err = ss.si.Create(l); err != nil {
		return pkg.NewAppError("could not create user", err, http.StatusInternalServerError)
	}
	return err
}

func (ss SignService) SignIn(l *model.LogIn) (token string, refresh string, err error) {
	if err := ss.va.Validate(l); err != nil {
		return "", "", fmt.Errorf("validate: %w", err)
	}
	DB, err := ss.si.Find(l.User)
	if err != nil {
		return "", "", pkg.NewAppError("email not found", err, http.StatusBadRequest)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(DB.Password), []byte(l.Password)); err != nil {
		return "", "", pkg.NewAppError("wrong password", err, http.StatusBadRequest)
	}
	token, err = ss.to.Create(DB.ID, uint(ss.va.UType()))
	if err != nil {
		return "", "", fmt.Errorf("save token: %w", err)
	}
	refresh, err = ss.createRefreshToken(DB.ID, pkg.UserType(ss.va.UType()))
	if err != nil {
		return "", "", fmt.Errorf("create refresh token: %w", err)
	}
	return
}

func (ss SignService) SignOut(refresh *string) error {
	jwt, err := ss.to.Validate(refresh)
	if err != nil {
		return fmt.Errorf("validate token: %w", err)
	}
	public := jwt.Public()
	id := public["id"].(float64)
	err = ss.re.Delete(uint(id))
	if err != nil {
		return pkg.NewAppError("could not delete token", err, http.StatusInternalServerError)
	}
	return nil
}

func (ss SignService) Refresh(refresh *string) (token string, err error) {
	uid, err := ss.validateRefreshToken(refresh)
	if err != nil {
		return "", pkg.NewAppError("failed validate refresh token", err, http.StatusBadRequest)
	}
	token, err = ss.to.Create(uid, uint(ss.va.UType()))
	if err != nil {
		return "", fmt.Errorf("create token: %w", err)
	}
	return
}

func (ss SignService) createRefreshToken(id uint, ut pkg.UserType) (string, error) {
	fgp, err := generateFgp(48)
	if err != nil {
		return "", pkg.NewAppError("could not create token", err, http.StatusInternalServerError)
	}
	refresh := model.Refresh{
		Hash:      fgp,
		ExpiresAt: time.Now().Add(168 * time.Hour), //TODO Change this number
		UserType:  int(ut),
	}
	err = ss.re.Create(&refresh)
	if err != nil {
		return "", pkg.NewAppError("fail to save token", err, http.StatusInternalServerError)
	}
	refreshToken, err := ss.to.CreateRefresh(refresh.ID, id, &fgp)
	if err != nil {
		return "", fmt.Errorf("create refresh token: %w", err)
	}
	return refreshToken, nil
}

func (ss SignService) validateRefreshToken(token *string) (uint, error) {
	jwt, err := ss.to.Validate(token)
	if err != nil {
		return 0, fmt.Errorf("validate token: %w", err)
	}
	public := jwt.Public()
	id := public["id"].(float64)
	re, err := ss.re.Find(uint(id))
	if err != nil {
		return 0, pkg.NewAppError("token not found", err, http.StatusBadRequest)
	}
	fgp := public["fgp"].(string)
	//hashFgp := hashFgp([]byte(fgp))
	if !strings.EqualFold(fgp, re.Hash) {
		return 0, pkg.NewAppError("wrong code", nil, http.StatusBadRequest)
	}
	uid, ok := public["uid"].(float64)
	if !ok {
		return 0, pkg.NewAppError("user id is not a number", err, http.StatusBadRequest)
	}
	return uint(uid), nil
}
