package controller

import (
	"errors"
	"fmt"
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

func isNotFoundErr(err error) bool {
	var nf interface{ IsNotFound() }
	return errors.As(err, &nf)
}

func (ss SignService) SignUp(l *model.LogIn) error {
	err := ss.va.Validate(l)
	if err != nil {
		return err
	}
	_, err = ss.si.Find(l.User)
	if err != nil && !isNotFoundErr(err) {
		return err
	}
	if err == nil {
		return ErrEmailAlreadyInUsed
	}
	pwd, err := hashAndSalt([]byte(l.Password))
	if err != nil {
		return err
	}
	l.Password = string(pwd)

	return ss.si.Create(l)
}

func (ss SignService) SignIn(l *model.LogIn) (token string, refresh string, err error) {
	if err := ss.va.Validate(l); err != nil {
		return "", "", err
	}
	DB, err := ss.si.Find(l.User)
	if err != nil {
		return "", "", fmt.Errorf("%w: %s", ErrUserNotFound, err)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(DB.Password), []byte(l.Password)); err != nil {
		return "", "", fmt.Errorf("%w: %s", ErrWrongPassword, err)
	}
	token, err = ss.to.Create(DB.ID, uint(ss.va.UType()))
	if err != nil {
		return "", "", err
	}
	refresh, err = ss.createRefreshToken(DB.ID, pkg.UserType(ss.va.UType()))
	return
}

func (ss SignService) SignOut(refresh *string) error {
	jwt, err := ss.to.Validate(refresh)
	if err != nil {
		return err
	}
	public := jwt.Public()
	id := public["id"].(float64)
	err = ss.re.Delete(uint(id))
	if err != nil {
		return fmt.Errorf("error at delete refresh token in DB: %w", err)
	}
	return nil
}

func (ss SignService) Refresh(refresh *string) (token string, err error) {
	uid, err := ss.validateRefreshToken(refresh)
	if err != nil {
		return "", fmt.Errorf("error at validateRefreshToken: %s", err)
	}
	token, err = ss.to.Create(uid, uint(ss.va.UType()))
	if err != nil {
		return "", fmt.Errorf("error at GenerateToken: %s", err)
	}
	return
}

func (ss SignService) createRefreshToken(id uint, ut pkg.UserType) (string, error) {
	fgp, err := generateFgp(48)
	if err != nil {
		return "", fmt.Errorf("error at GenerateFgp: %s", err)
	}
	refresh := model.Refresh{
		Hash:      fgp,
		ExpiresAt: time.Now().Add(168 * time.Hour),
		UserType:  int(ut),
	}
	err = ss.re.Create(&refresh)
	if err != nil {
		return "", err
	}
	refreshToken, err := ss.to.CreateRefresh(refresh.ID, id, &fgp)
	if err != nil {
		return "", fmt.Errorf("error at GenerateRefreshToken: %s", err)
	}
	return refreshToken, nil
}

func (ss SignService) validateRefreshToken(token *string) (uint, error) {
	jwt, err := ss.to.Validate(token)
	if err != nil {
		return 0, err
	}
	public := jwt.Public()
	id := public["id"].(float64)
	re, err := ss.re.Find(uint(id))
	if err != nil {
		return 0, err
	}
	fgp := public["fgp"].(string)
	//hashFgp := hashFgp([]byte(fgp))
	if !strings.EqualFold(fgp, re.Hash) {
		return 0, ErrInvalidRefreshToken
	}
	uid := public["uid"].(float64)
	return uint(uid), nil
}
