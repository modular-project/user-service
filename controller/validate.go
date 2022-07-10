package controller

import (
	"net/http"
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
		return pkg.NewAppError("invalid email", nil, http.StatusBadRequest)
	}
	if ok := isPasswordValid(l.Password); !ok {
		return pkg.NewAppError("invalid password", nil, http.StatusBadRequest)
	}
	return nil
}

func (uv userValidate) UType() pkg.UserType {
	return uv.ut
}
