package controller_test

import (
	"errors"
	"net/http"
	"testing"
	"users-service/controller"
	"users-service/mocks"
	"users-service/model"
	"users-service/pkg"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGenerateCode(t *testing.T) {
	ve := mocks.NewVerificationStorager(t)
	u := mocks.NewUserStorager(t)
	mail := mocks.NewMailer(t)

	ve.On("Create", mock.Anything).Return(nil)
	mail.On("Confirm", "mail@mail.com", mock.Anything).Return(nil)
	tests := []struct {
		give     uint
		wantCode int
		wantErr  error
	}{
		{1, 0, nil},
		{2, http.StatusBadRequest, errors.New("not found")},
	}

	for _, tt := range tests {
		u.On("Find", tt.give).Return(model.User{Email: "mail@mail.com"}, tt.wantErr)
		us := controller.NewUserService(u, ve, mail)
		gotErr := us.GenerateCode(tt.give)
		gotCode, _ := pkg.FindError(gotErr)
		assert.Equal(t, tt.wantCode, gotCode)
	}
}
