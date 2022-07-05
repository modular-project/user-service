package controller_test

import (
	"errors"
	"log"
	"testing"
	"users-service/authorization"
	"users-service/controller"
	"users-service/model"
	"users-service/storage"
)

var TestConfigDB storage.DBConnection = storage.DBConnection{
	TypeDB:   storage.POSTGRESQL,
	User:     "admin_restaurant",
	Password: "RestAuraNt_pgsql.561965697",
	Host:     "localhost",
	Port:     "5433",
	NameDB:   "testing",
}

func dropsTables(t *testing.T, tables ...interface{}) {
	err := storage.Drop(tables...)
	if err != nil {
		t.Fatalf("Failed to clean database: %s", err)
	}
}

func TestSingUpIntegration(t *testing.T) {
	err := storage.NewDB(TestConfigDB)
	if err != nil {
		t.Fatalf("NewGormDB: %s", err)
	}
	uc := controller.NewSignService(storage.NewRefreshStore(), controller.NewUserValidate(), storage.NewUserSignStore(), nil)
	models := []interface{}{model.User{}, model.Role{}, model.UserRole{}}
	err = storage.Migrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })
	//insertUsers(uc)
	testCase := []struct {
		in   model.LogIn
		want error
	}{
		{model.LogIn{
			User:     "cualquiera",
			Password: "Password12345.",
		}, controller.ErrEmailNotValid},
		{model.LogIn{
			User:     "usuario@mail.com",
			Password: "as",
		}, controller.ErrPasswordNotValid},
		{model.LogIn{
			User:     "usuario@mail.com",
			Password: "Password12345.",
		}, nil},
		{model.LogIn{
			User:     "usuario@mail.com",
			Password: "Password12345.",
		}, controller.ErrEmailAlreadyInUsed},
		{model.LogIn{
			User:     "usuario1@mail.com",
			Password: " ",
		}, controller.ErrPasswordNotValid},
		{model.LogIn{
			User:     "usuario1@mail.com",
			Password: "Password12345. ",
		}, nil},
	}

	for _, tc := range testCase {
		log.Printf("%+v", tc.in)
		err := uc.SignUp(&tc.in)
		if err != tc.want {
			t.Errorf("Got: %s, want: %s", err, tc.want)
		}
	}
}

func TestSingInIntegration(t *testing.T) {
	err := storage.NewDB(TestConfigDB)
	if err != nil {
		t.Fatalf("NewGormDB: %s", err)
	}
	authorization.LoadCertificates(authorization.RSA512)
	uc := controller.NewSignService(storage.NewRefreshStore(), controller.NewUserValidate(), storage.NewUserSignStore(), authorization.NewToken())
	models := []interface{}{model.User{}, model.Role{}, model.UserRole{}, model.Refresh{}}
	err = storage.Migrate(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	t.Cleanup(func() { dropsTables(t, models...) })

	// err = authorization.LoadCertificates()
	if err != nil {
		t.Fatalf("error at loadCertificates: %s", err)
	}
	insertUsers(uc)
	testCase := []struct {
		in   model.LogIn
		want error
	}{
		{
			model.LogIn{User: "valid@account.com", Password: "PassOk1234. "}, nil,
		},
		{
			model.LogIn{User: "notAnEmail", Password: "ValidPass1234."}, controller.ErrEmailNotValid,
		},
		{
			model.LogIn{User: "NotAnCreatedAccount@mail.com", Password: "ValidPass123."}, controller.ErrUserNotFound,
		},
		{
			model.LogIn{User: "valid@account.com", Password: "WrongPass123."}, controller.ErrWrongPassword,
		},
		{
			model.LogIn{User: "valid@account.com", Password: "notapass"}, controller.ErrPasswordNotValid,
		},
	}

	for _, tc := range testCase {
		_, _, gotErr := uc.SignIn(&tc.in)
		if !errors.Is(gotErr, tc.want) {
			t.Logf("%+v", tc)
			t.Errorf("Got: %s, want: %s", gotErr, tc.want)
		}
	}
}

// func TestVerifyUser(t *testing.T) {
// 	err := storage.NewDB(TestConfigDB)
// 	if err != nil {
// 		t.Fatalf("NewGormDB: %s", err)
// 	}
// 	authorization.LoadCertificates(authorization.RSA512)
// 	us := controller.NewSignService(storage.NewRefreshStore(), controller.NewUserValidate(), storage.NewUserSignStore(), authorization.NewToken())
// 	uc := controller.NewUserService()
// 	models := []interface{}{model.User{}, model.Role{}, model.UserRole{}, model.Verification{}}
// 	err = storage.Migrate(models...)
// 	if err != nil {
// 		t.Fatalf("Failed to Create tables: %s", err)
// 	}
// 	t.Cleanup(func() { dropsTables(t, models...) })
// 	testCase := []struct {
// 		userID uint
// 		code   string
// 		err    error
// 	}{
// 		{1, "AA3AAA", nil},
// 		{2, "", controller.ErrNullCode},
// 		{3, "NOCODE", controller.ErrInvalidCode},
// 		{4, "AA3AAA", controller.ErrExpiredCode},
// 		{5, "AA3AAA", controller.ErrNoRowsAffected},
// 		{6, "AA3AAA", controller.ErrCodeNotFound},
// 	}
// 	//user := model.LogIn{User: "a", Password: "a"}
// 	insertUsers(us)
// 	//err = storage.DB().Create(&user).Error
// 	verifications := []model.Verification{
// 		{UserID: 1, Code: "AA3AAA", ExpiresAt: time.Now().Add(time.Minute * 15)},
// 		{UserID: 2, Code: "AA3AAA", ExpiresAt: time.Now().Add(time.Minute * 15)},
// 		{UserID: 3, Code: "AA3AAA", ExpiresAt: time.Now().Add(time.Minute * 15)},
// 		{UserID: 4, Code: "AA3AAA", ExpiresAt: time.Now().Add(-time.Minute * 15)},
// 		{UserID: 5, Code: "AA3AAA", ExpiresAt: time.Now().Add(time.Minute * 15)},
// 	}
// 	for _, v := range verifications {
// 		err := uc.Ver.Create(&v)
// 		if err != nil {
// 			t.Logf("%+v", v)
// 			t.Fatalf("error at create verifications: %s", err)
// 		}
// 	}

// 	for _, tc := range testCase {
// 		err := uc.Verify(tc.userID, tc.code)
// 		if !errors.Is(tc.err, err) {
// 			t.Logf("%+v", tc)
// 			t.Errorf("got error: %s, want error: %s", err, tc.err)
// 		}
// 	}
// }

// func TestUpdateUserData(t *testing.T) {
// 	err := storage.NewDB(TestConfigDB)
// 	if err != nil {
// 		t.Fatalf("NewGormDB: %s", err)
// 	}
// 	us := controller.NewSignService(storage.NewRefreshStore(), controller.NewUserValidate(), storage.NewUserSignStore())
// 	uc := controller.NewUserService()
// 	models := []interface{}{model.LogIn{}, model.UserRole{}}
// 	err = storage.GormMigrate(models...)
// 	if err != nil {
// 		t.Fatalf("Failed to Create tables: %s", err)
// 	}
// 	t.Cleanup(func() { dropsTables(t, models...) })
// 	now := time.Now()
// 	before := now.Add(-time.Hour * 5)
// 	testCase := []struct {
// 		in   model.User
// 		name string
// 		url  string
// 		err  error
// 	}{
// 		{model.User{Model: model.Model{ID: 1, UpdatedAt: before}}, "", "", nil},
// 		{model.User{Model: model.Model{ID: 1, CreatedAt: before}, BirthDate: time.Date(2000, time.October, 29, 0, 0, 0, 0, time.UTC)}, "nicolas", "foto/nicolas.jpg", nil},
// 		{model.User{Model: model.Model{ID: 6}, BirthDate: time.Time{}}, "", "", controller.ErrNoRowsAffected},
// 		{model.User{Model: model.Model{ID: 1}, BirthDate: time.Time{}}, "nil", "nuevafoto/nicolas", nil},
// 		{model.User{Model: model.Model{ID: 2, UpdatedAt: before}}, "false info", "falseinfo.jpg", controller.ErrUnauthorizedUser},
// 	}
// 	oldUsers := []model.LogIn{
// 		{User: "valid@account.com", Password: "PassOk1234. "},
// 		{User: "valid2@account.com", Password: "PassOk1234. "},
// 	}
// 	for _, ou := range oldUsers {
// 		res := us.SS.Create(&ou)
// 		if err != nil {
// 			t.Fatalf("fatal at create user: %s", res.Error)
// 		}
// 	}
// 	storage.DB().Create(&model.UserRole{UserID: 2})
// 	for i, tc := range testCase {
// 		if tc.name != "nil" {
// 			tc.in.Name = &tc.name
// 		}
// 		if tc.url != "nil" {
// 			tc.in.URL = &tc.url
// 		}
// 		err := controller.UpdateUserData(&tc.in)
// 		if !errors.Is(err, tc.err) {
// 			t.Logf("%+v", tc.in)
// 			t.Errorf("got error: %s, want error: %s", err, tc.err)
// 		}
// 		if err == nil {
// 			u := model.LogIn{}
// 			storage.DB().Where("id = ?", tc.in.ID).First(&u)
// 			t.Logf("%+v", u)
// 			if u.UpdatedAt.Equal(tc.in.UpdatedAt) || !u.BirthDate.Equal(tc.in.BirthDate) || !reflect.DeepEqual(u.Name, tc.in.Name) || !reflect.DeepEqual(u.URL, tc.in.URL) {
// 				t.Errorf("%d - got user: Name: %+v, Url: %+v, Birth: %s, UpdatedAt: %s CreatedAt: %s, want user: Name: %s, Url: %s, Birth: %s, UpdatedAt: %s CreatedAt: %s",
// 					i, u.Name, u.URL, u.BirthDate.UTC(), u.UpdatedAt.UTC(), u.CreatedAt.UTC(), tc.name, tc.url, tc.in.BirthDate.UTC(), tc.in.UpdatedAt.UTC(), tc.in.CreatedAt.UTC())
// 			}
// 		}
// 	}

// }

// func TestChangePassword(t *testing.T) {
// 	storage.New(storage.TESTING)
// 	models := []interface{}{model.LogIn{}}
// 	err := storage.DB().AutoMigrate(models...)
// 	if err != nil {
// 		t.Fatalf("Failed to Create tables: %s", err)
// 	}
// 	t.Cleanup(func() { dropsTables(t, models...) })
// 	testCase := []struct {
// 		in  model.LogIn
// 		err error
// 	}{
// 		{model.LogIn{Model: model.Model{ID: 1}, Password: "NewPassword12345."}, nil},
// 		{model.LogIn{Model: model.Model{ID: 1}, Password: "invalidPassword"}, controller.ErrPasswordNotValid},
// 		{model.LogIn{Model: model.Model{ID: 7}, Password: "NewPassword12345."}, controller.ErrNoRowsAffected},
// 	}
// 	user := model.LogIn{Password: "NewPassword12345."}
// 	storage.DB().Create(&user)
// 	for i, tc := range testCase {
// 		err := controller.ChangeUserPassword(tc.in.ID, &tc.in.Password)
// 		if !errors.Is(err, tc.err) {
// 			t.Logf("%d - %+v", i, tc.in)
// 			t.Errorf("got error: %s, want error: %s", err, tc.err)
// 		}
// 	}

// }

func insertUsers(uc controller.SignService) {
	data := []model.LogIn{
		{User: "valid@account.com", Password: "PassOk1234. "},
	}
	for _, d := range data {
		uc.SignUp(&d)
	}
}
