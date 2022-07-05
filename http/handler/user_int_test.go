package handler_test

// import (
// 	"bytes"
// 	"encoding/json"
// 	"errors"
// 	"net/http"
// 	"net/http/httptest"
// 	"testing"
// 	"users-service/authorization"
// 	"users-service/controller"
// 	"users-service/http/handler"
// 	"users-service/model"
// 	"users-service/storage"

// 	"github.com/labstack/echo"
// )

// type ResponseJson struct {
// 	Message string `json:"msg"`
// }

// type TestCase struct {
// 	in           model.User
// 	responseCode int
// 	message      error
// }

// func dropsTables(t *testing.T, tables ...interface{}) {
// 	err := storage.DB().Migrator().DropTable(tables...)
// 	if err != nil {
// 		t.Fatalf("Failed to clean database: %s", err)
// 	}
// }

// func TestSingUpIntegration(t *testing.T) {
// 	storage.New(storage.TESTING)
// 	models := []interface{}{model.User{}, model.Role{}, model.UserRole{}}
// 	err := storage.DB().AutoMigrate(models...)
// 	if err != nil {
// 		t.Fatalf("Failed to Create tables: %s", err)
// 	}
// 	t.Cleanup(func() { dropsTables(t, models...) })

// 	testCase := []TestCase{
// 		{model.User{
// 			Email:    "cualquiera",
// 			Password: "Password12345.",
// 		}, http.StatusBadRequest, controller.ErrEmailNotValid},
// 		{model.User{
// 			Email:    "usuario@mail.com",
// 			Password: "as",
// 		}, http.StatusBadRequest, controller.ErrPasswordNotValid},
// 		{model.User{
// 			Email:    "usuario@mail.com",
// 			Password: "Password12345.",
// 		}, http.StatusCreated, nil},
// 		{model.User{
// 			Email:    "usuario@mail.com",
// 			Password: "Password12345.",
// 		}, http.StatusBadRequest, controller.ErrEmailAlreadyInUsed},
// 		{model.User{
// 			Email:    "usuario1@mail.com",
// 			Password: " ",
// 		}, http.StatusBadRequest, controller.ErrPasswordNotValid},
// 		{model.User{
// 			Email:    "usuario1@mail.com",
// 			Password: "Password12345. ",
// 		}, http.StatusCreated, nil},
// 	}

// 	for _, tc := range testCase {
// 		data, err := json.Marshal(tc.in)
// 		if err != nil {
// 			t.Fatalf("Error at marshal: %s", err)
// 		}

// 		r := httptest.NewRequest(http.MethodPost, "/signup/", bytes.NewBuffer(data))
// 		r.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		e := echo.New()
// 		ctx := e.NewContext(r, w)
// 		err = handler.SignUp(ctx)
// 		if err != nil {
// 			t.Errorf("Got error at SignUp: %s", err)
// 		}

// 		if ctx.Response().Status != tc.responseCode {
// 			t.Errorf("Got status code: %d, want: %d", ctx.Response().Status, tc.responseCode)
// 		}
// 		if ctx.Response().Header().Get("Content-Type") == "application/json; charset=UTF-8" {
// 			m := ResponseJson{}
// 			err = json.NewDecoder(w.Body).Decode(&m)
// 			if err != nil {
// 				t.Errorf("Got error: %s at decode", err)
// 			}

// 			if errors.Is(tc.message, errors.New(m.Message)) {
// 				t.Errorf("Got message: %s, want: %s", m.Message, tc.message)
// 			}
// 		}

// 	}
// }

// func TestSingInIntegration(t *testing.T) {
// 	storage.New(storage.TESTING)
// 	authorization.LoadCertificates()
// 	models := []interface{}{model.User{}, model.Role{}, model.UserRole{}, model.Refresh{}}
// 	err := storage.DB().AutoMigrate(models...)
// 	if err != nil {
// 		t.Fatalf("Failed to Create tables: %s", err)
// 	}
// 	t.Cleanup(func() { dropsTables(t, models...) })
// 	insertUsers()
// 	testCase := []struct {
// 		in           model.User
// 		responseCode int
// 		message      error
// 	}{
// 		{
// 			model.User{Email: "valid@account.com", Password: "PassOk1234."}, http.StatusOK, nil,
// 		},
// 		{
// 			model.User{Email: "notAnEmail", Password: "ValidPass1234."}, http.StatusBadRequest, controller.ErrEmailNotValid,
// 		},
// 		{
// 			model.User{Email: "NotAnCreatedAccount@mail.com", Password: "ValidPass123."}, http.StatusBadRequest, controller.ErrUserNotFound,
// 		},
// 		{
// 			model.User{Email: "valid@account.com", Password: "WrongPass123."}, http.StatusBadRequest, controller.ErrWrongPassword,
// 		},
// 		{
// 			model.User{Email: "valid@account.com", Password: "notapass"}, http.StatusBadRequest, controller.ErrPasswordNotValid,
// 		},
// 	}

// 	for _, tc := range testCase {
// 		data, err := json.Marshal(tc.in)
// 		if err != nil {
// 			t.Fatalf("Error at marshal: %s", err)
// 		}

// 		r := httptest.NewRequest(http.MethodPost, "/signin/", bytes.NewBuffer(data))
// 		r.Header.Set("Content-Type", "application/json")
// 		w := httptest.NewRecorder()
// 		e := echo.New()
// 		c := e.NewContext(r, w)
// 		err = handler.SignIn(c)
// 		if err != nil {
// 			t.Errorf("error at SignIn: %s", err)
// 		}
// 		if c.Response().Status != tc.responseCode {
// 			t.Errorf("Got status code: %d, want: %d", c.Response().Status, tc.responseCode)
// 		}
// 		if c.Response().Header().Get("Content-Type") == "application/json; charset=UTF-8" {
// 			m := ResponseJson{}
// 			err = json.NewDecoder(w.Body).Decode(&m)
// 			if err != nil {
// 				t.Errorf("Got error: %s at decode", err)
// 			}
// 			if c.Response().Status != http.StatusOK {
// 				if errors.Is(tc.message, errors.New(m.Message)) {
// 					t.Errorf("Got message: %s, want: %s", m.Message, tc.message)
// 				}
// 			} else {
// 				_, err = authorization.ValidateToken(&m.Message)
// 				t.Log(m.Message)
// 				if err != nil {
// 					t.Errorf("Error at validateToken: %s", err)
// 				}
// 				cookie := w.Result().Cookies()[0]
// 				r = httptest.NewRequest(http.MethodGet, "/api/v1/user/refresh/", nil)
// 				w2 := httptest.NewRecorder()
// 				c2 := e.NewContext(r, w2)
// 				r.AddCookie(cookie)
// 				testRefresh(t, c2, w2)
// 			}
// 		}
// 	}
// }

// func testRefresh(t *testing.T, c echo.Context, w *httptest.ResponseRecorder) {
// 	err := handler.Refresh(c)
// 	if err != nil {
// 		t.Fatalf("Error at refresh: %s", err)
// 	}
// 	m := ResponseJson{}
// 	err = json.NewDecoder(w.Body).Decode(&m)
// 	if err != nil {
// 		t.Errorf("Got error: %s at decode", err)
// 	}
// 	_, err = authorization.ValidateToken(&m.Message)
// 	t.Log(m.Message)
// 	if err != nil {
// 		t.Errorf("Error at validateToken: %s", err)
// 	}
// }

// func insertUsers() {
// 	data := []model.User{
// 		{Email: "valid@account.com", Password: "PassOk1234."},
// 	}
// 	for _, d := range data {
// 		controller.SignUp(&d)
// 	}
// }
