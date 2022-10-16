package handler

import (
	"fmt"
	"net/http"
	"users-service/model"
	"users-service/pkg"

	"github.com/labstack/echo"
)

type SignServicer interface {
	SignUp(*model.LogIn) error
	SignIn(*model.LogIn) (token string, refresh string, err error)
	SignOut(refresh *string) error
	Refresh(*string) (token string, err error)
}

type UserServicer interface {
	Data(uint) (model.User, error)
	Verify(uint, string) error
	GenerateCode(uint) error
	UpdateData(*model.User) error
	ChangePassword(uint, *string) error
}

type UserUC struct {
	us UserServicer
	ss SignServicer
}

func NewUserUC(uc UserServicer, ss SignServicer) UserUC {
	return UserUC{uc, ss}
}

func (uuc UserUC) SignUp(c echo.Context) error {
	var err error
	m := &model.LogIn{}
	if err = c.Bind(m); err != nil {
		return pkg.NewAppError("Fail at bind Log In", err, http.StatusBadRequest)
	}
	err = uuc.ss.SignUp(m)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusCreated)
}

func (uuc UserUC) SignIn(c echo.Context) error {
	var err error
	m := &model.LogIn{}
	if err = c.Bind(m); err != nil {
		return pkg.NewAppError("Fail at bind Log In", err, http.StatusBadRequest)
	}
	token, refresh, err := uuc.ss.SignIn(m)
	if err != nil {
		return err
	}
	createRefreshCookie(c, refresh, "/api/v1/user/refresh/")
	return c.JSON(http.StatusOK, createResponse(token))
}

func (uuc UserUC) Refresh(c echo.Context) error {
	fgp, err := c.Cookie("refresh")
	if err != nil {
		return pkg.NewAppError("Fail at get cookie", err, http.StatusBadRequest)
	}
	token, err := uuc.ss.Refresh(&fgp.Value)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, createResponse(token))
}

func (uuc UserUC) GetUserData(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	m, err := uuc.us.Data(userID)
	if err != nil {
		return fmt.Errorf("GetUserData: %s", err)
	}
	m.Password = ""
	return c.JSON(http.StatusOK, &m)
}

// Update bassc user data like name, brithdate and image
func (uuc UserUC) UpdateUserData(c echo.Context) error {
	m := &model.User{}
	if err := c.Bind(m); err != nil {
		return pkg.NewAppError("Fail at bind user", err, http.StatusBadRequest)
	}
	id, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	m.ID = id
	err = uuc.us.UpdateData(m)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (uuc UserUC) SignOut(c echo.Context) error {
	fgp, err := c.Cookie("refresh")
	if err != nil {
		return pkg.NewAppError("Fail at get cookie", err, http.StatusBadRequest)
	}
	err = uuc.ss.SignOut(&fgp.Value)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (uuc UserUC) VerifyUser(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	code := c.QueryParam("code")
	err = uuc.us.Verify(userID, code)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (uuc UserUC) GenerateVerificationCode(c echo.Context) error {
	userID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	err = uuc.us.GenerateCode(userID)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (uuc UserUC) ChangePassword(c echo.Context) error {
	m := &model.User{}
	if err := c.Bind(m); err != nil {
		return pkg.NewAppError("Fail at bind user", err, http.StatusBadRequest)
	}
	id, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	err = uuc.us.ChangePassword(id, &m.Password)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
