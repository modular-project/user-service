package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"users-service/model"
	"users-service/pkg"

	"github.com/labstack/echo"
)

type KitchenServicer interface {
	GetInESTB(uint) ([]model.Kitchen, error)
	Delete(eID, kID uint) error
	Update(eID, kID uint, k *model.Kitchen) error
}

type KitchenUC struct {
	ss SignServicer
	ks KitchenServicer
}

func NewKitchenUC(ss SignServicer, ks KitchenServicer) KitchenUC {
	return KitchenUC{ss: ss, ks: ks}
}

func (kuc KitchenUC) SignUp(c echo.Context) error {
	var login model.LogIn
	if err := c.Bind(&login); err != nil {
		return pkg.NewAppError("failed at bind login", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	login.ID = ur.EstablishmentID
	if err := kuc.ss.SignUp(&login); err != nil {
		return fmt.Errorf("ss.SignUp: %w", err)
	}
	return c.NoContent(http.StatusOK)
}

func (kuc KitchenUC) SignIn(c echo.Context) error {
	var login model.LogIn
	if err := c.Bind(&login); err != nil {
		return pkg.NewAppError("failed at bind login", err, http.StatusBadRequest)
	}
	t, tr, err := kuc.ss.SignIn(&login)
	if err != nil {
		return fmt.Errorf("ss.SignUp: %w", err)
	}
	createRefreshCookie(c, tr, "/api/v1/kitchen/refresh/")
	return c.JSON(http.StatusOK, createResponse(t))
}

func (kuc KitchenUC) SignOut(c echo.Context) error {
	fgp, err := c.Cookie("refresh")
	if err != nil {
		return pkg.NewAppError("Fail at get cookie", err, http.StatusBadRequest)
	}
	err = kuc.ss.SignOut(&fgp.Value)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (kuc KitchenUC) Refresh(c echo.Context) error {
	fgp, err := c.Cookie("refresh")
	if err != nil {
		return pkg.NewAppError("Fail at get cookie", err, http.StatusBadRequest)
	}
	token, err := kuc.ss.Refresh(&fgp.Value)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, createResponse(token))
}

func (kuc KitchenUC) GetInESTB(c echo.Context) error {
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	kits, err := kuc.ks.GetInESTB(ur.EstablishmentID)
	if err != nil {
		return err
	}
	if kits == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, kits)
}

func (kuc KitchenUC) Update(c echo.Context) error {
	var kit model.Kitchen
	if err := c.Bind(&kit); err != nil {
		return pkg.NewAppError("failed at bind login", err, http.StatusBadRequest)
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("failed to get id from path param", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	if err := kuc.ks.Update(ur.EstablishmentID, uint(id), &kit); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (kuc KitchenUC) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("failed to get id from path param", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	if err := kuc.ks.Delete(ur.EstablishmentID, uint(id)); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
