package controller_test

import (
	"errors"
	"testing"
	"users-service/controller"
	"users-service/mocks"
	"users-service/model"
	"users-service/pkg"

	"github.com/stretchr/testify/mock"
)

func TestGenerateCode(t *testing.T) {
	ve := mocks.NewVerificationStorager(t)
	u := mocks.NewUserStorager(t)
	mail := mocks.NewMailer(t)

	ve.On("Create", mock.Anything).Return(nil)
	mail.On("Confirm", "mail@mail.com", mock.Anything).Return(nil)
	tests := []struct {
		give    uint
		wantErr error
	}{
		{1, nil},
		{2, pkg.ErrNoRowsAffected},
	}

	for _, tt := range tests {
		u.On("Find", tt.give).Return(model.User{Email: "mail@mail.com"}, tt.wantErr)
		us := controller.NewUserService(u, ve, mail)
		gotErr := us.GenerateCode(tt.give)
		if !errors.Is(gotErr, tt.wantErr) {
			t.Errorf("got err: %s, want error: %s", gotErr, tt.wantErr)
		}
	}
}
