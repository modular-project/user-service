package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"users-service/http/handler"
	"users-service/http/middleware"
	"users-service/mocks"
	"users-service/model"
	"users-service/pkg"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type ResponseJson struct {
	Message string `json:"msg"`
}

type TestCase struct {
	give     model.LogIn
	wantCode int
}

func TestSingUp(t *testing.T) {
	tests := []TestCase{
		{model.LogIn{
			User:     "cualquiera",
			Password: "Password12345.",
		}, http.StatusBadRequest},
		{model.LogIn{
			User:     "usuario@mail.com",
			Password: "as",
		}, http.StatusBadRequest},
		{model.LogIn{
			User:     "usuario@mail.com",
			Password: "Password12345.",
		}, http.StatusCreated},
		{model.LogIn{
			User:     "usuario@mail.com",
			Password: "Password12345.",
		}, http.StatusBadRequest},
		{model.LogIn{
			User:     "usuario1@mail.com",
			Password: " ",
		}, http.StatusBadRequest},
		{model.LogIn{
			User:     "usuario1@mail.com",
			Password: "Password12345. ",
		}, http.StatusCreated},
	}
	for _, tt := range tests {
		si := mocks.NewSignUCer(t)
		assert := assert.New(t)
		if tt.wantCode != http.StatusCreated {
			si.On("SignUp", mock.Anything).Return(pkg.NewAppError("Fail at sing up", nil, http.StatusBadRequest)).Once()
		} else {
			si.On("SignUp", mock.Anything).Return(nil).Once()
		}
		data, err := json.Marshal(tt.give)
		if err != nil {
			t.Fatalf("Error at marshal: %s", err)
		}

		r := httptest.NewRequest(http.MethodPost, "/signup/", bytes.NewBuffer(data))
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		e := echo.New()
		m := middleware.NewMiddleware(nil)
		ctx := e.NewContext(r, w)
		h := handler.NewUserUC(nil, si)
		gotErr := h.SignUp(ctx)
		if gotErr != nil {
			m.Errors(gotErr, ctx)
		}

		assert.Equal(tt.wantCode, ctx.Response().Status)

	}
}

// func TestSingIn(t *testing.T) {
// 	storage.New(storage.TESTING)
// 	authorization.LoadCertificates()
// 	models := []interface{}{model.LogIn{}, model.Role{}, model.UserRole{}, model.Refresh{}}
// 	err := storage.DB().AutoMigrate(models...)
// 	if err != nil {
// 		t.Fatalf("Failed to Create tables: %s", err)
// 	}
// 	t.Cleanup(func() { dropsTables(t, models...) })
// 	insertUsers()
// 	testCase := []struct {
// 		in           model.LogIn
// 		responseCode int
// 		message      error
// 	}{
// 		{
// 			model.LogIn{Email: "valid@account.com", Password: "PassOk1234."}, http.StatusOK, nil,
// 		},
// 		{
// 			model.LogIn{Email: "notAnEmail", Password: "ValidPass1234."}, http.StatusBadRequest, controller.ErrEmailNotValid,
// 		},
// 		{
// 			model.LogIn{Email: "NotAnCreatedAccount@mail.com", Password: "ValidPass123."}, http.StatusBadRequest, controller.ErrUserNotFound,
// 		},
// 		{
// 			model.LogIn{Email: "valid@account.com", Password: "WrongPass123."}, http.StatusBadRequest, controller.ErrWrongPassword,
// 		},
// 		{
// 			model.LogIn{Email: "valid@account.com", Password: "notapass"}, http.StatusBadRequest, controller.ErrPasswordNotValid,
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

// // func testRefresh(t *testing.T, c echo.Context, w *httptest.ResponseRecorder) {
// // 	err := handler.Refresh(c)
// // 	if err != nil {
// // 		t.Fatalf("Error at refresh: %s", err)
// // 	}
// // 	m := ResponseJson{}
// // 	err = json.NewDecoder(w.Body).Decode(&m)
// // 	if err != nil {
// // 		t.Errorf("Got error: %s at decode", err)
// // 	}
// // 	_, err = authorization.ValidateToken(&m.Message)
// // 	t.Log(m.Message)
// // 	if err != nil {
// // 		t.Errorf("Error at validateToken: %s", err)
// // 	}
// // }

// // func insertUsers() {
// // 	data := []model.LogIn{
// // 		{Email: "valid@account.com", Password: "PassOk1234."},
// // 	}
// // 	for _, d := range data {
// // 		controller.SignUp(&d)
// // 	}
// // }
