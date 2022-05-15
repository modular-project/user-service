package controller_test

import (
	"testing"
	"users-service/controller"
	"users-service/model"
	"users-service/storage"
)

func dropsTables(t *testing.T, tables ...interface{}) {
	err := storage.DB().Migrator().DropTable(tables...)
	if err != nil {
		t.Fatalf("Failed to clean database: %s", err)
	}
}
func TestSingUpIntegration(t *testing.T) {
	storage.New(storage.TESTING)
	models := []interface{}{model.Account{}, model.Role{}, model.AccountRole{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })

	testCase := []struct {
		in   model.Account
		want error
	}{
		{model.Account{
			Email:    "cualquiera",
			Password: "Password12345.",
		}, controller.ErrEmailNotValid},
		{model.Account{
			Email:    "usuario@mail.com",
			Password: "as",
		}, controller.ErrPasswordNotValid},
		{model.Account{
			Email:    "usuario@mail.com",
			Password: "Password12345.",
		}, nil},
		{model.Account{
			Email:    "usuario@mail.com",
			Password: "Password12345.",
		}, controller.ErrEmailAlreadyInUsed},
		{model.Account{
			Email:    "usuario1@mail.com",
			Password: " ",
		}, controller.ErrPasswordNotValid},
		{model.Account{
			Email:    "usuario1@mail.com",
			Password: "Password12345. ",
		}, nil},
	}

	for _, tc := range testCase {
		err := controller.SignUp(&tc.in)
		if err != tc.want {
			t.Errorf("Got: %s, want: %s", err, tc.want)
		}
	}
}

func TestSingInIntegration(t *testing.T) {
	storage.New(storage.TESTING)
	models := []interface{}{model.Account{}, model.Role{}, model.AccountRole{}}
	err := storage.DB().AutoMigrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	insertUsers()
	testCase := []struct {
		in   model.Account
		want error
	}{
		{
			model.Account{Email: "valid@account.com", Password: "PassOk1234. "}, nil,
		},
		{
			model.Account{Email: "notAnEmail", Password: "ValidPass1234."}, controller.ErrEmailNotValid,
		},
		{
			model.Account{Email: "NotAnCreatedAccount@mail.com", Password: "ValidPass123."}, controller.ErrUserNotFound,
		},
		{
			model.Account{Email: "valid@account.com", Password: "WrongPass123."}, controller.ErrWrongPassword,
		},
		{
			model.Account{Email: "valid@account.com", Password: "notapass"}, controller.ErrPasswordNotValid,
		},
	}

	for _, tc := range testCase {
		_, gotErr := controller.SignIn(&tc.in)
		if gotErr != tc.want {
			t.Errorf("Got: %s, want: %s", gotErr, tc.want)
		}
	}
}

func insertUsers() {
	data := []model.Account{
		{Email: "valid@account.com", Password: "PassOk1234. "},
	}
	for _, d := range data {
		controller.SignUp(&d)
	}
}
