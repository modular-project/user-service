package controller

import (
	"net/http"
	"users-service/model"
	"users-service/pkg"
)

type userValidate struct {
	ut pkg.UserType
}

type kitchenValidate struct {
	ut pkg.UserType
}

func NewKitchenValidate() kitchenValidate {
	return kitchenValidate{pkg.KITCHEN}
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

func (kv kitchenValidate) isUserValid(u string) bool {
	return len(u) >= 5
}

func (kv kitchenValidate) Validate(l *model.LogIn) error {
	if !kv.isUserValid(l.User) {
		return pkg.NewAppError("invalid username", nil, http.StatusBadRequest)
	}
	if !isPasswordValid(l.Password) {
		return pkg.NewAppError("invalid password", nil, http.StatusBadRequest)
	}
	return nil
}

func (kv kitchenValidate) UType() pkg.UserType {
	return kv.ut
}
