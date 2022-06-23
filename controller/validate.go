package controller

import (
	"users-service/model"
	"users-service/pkg"
)

type userValidate struct {
	ut pkg.UserType
}

func NewUserValidate() userValidate {
	return userValidate{pkg.USER}
}

func (uv userValidate) Validate(l *model.LogIn) error {
	if ok := isEmailValid(l.User); !ok {
		return ErrEmailNotValid
	}
	if ok := isPasswordValid(l.Password); !ok {
		return ErrPasswordNotValid
	}
	return nil
}

func (uv userValidate) UType() pkg.UserType {
	return uv.ut
}
