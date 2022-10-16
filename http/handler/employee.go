package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"users-service/model"
	"users-service/pkg"

	"github.com/labstack/echo"
)

type EMPLService interface {
	Update(from model.UserRole, target uint, u *model.User) error
	Self(uint) (model.UserJobs, error)
	Get(from model.UserRole, target uint) (model.UserJobs, error)
	Search(*model.SearchEMPL) ([]model.User, error)
	SearchWaiters(uint, *model.Search) ([]model.User, error)
	Hire(model.UserRole, string, *model.UserRole) error
	HireWaiter(model.UserRole, string, float64) error
	Fire(from model.UserRole, target uint) error
}

type EMPLuc struct {
	es EMPLService
}

func NewEMPLUC(es EMPLService) EMPLuc {
	return EMPLuc{es: es}
}

func (eu EMPLuc) Update(c echo.Context) error {
	u := &model.User{}
	if err := c.Bind(u); err != nil {
		return pkg.NewAppError("Fail at bind user", err, http.StatusBadRequest)
	}
	tID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	f, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	err = eu.es.Update(f, uint(tID), u)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Search(c echo.Context) error {
	s := model.SearchEMPL{}
	if err := c.Bind(&s); err != nil {
		return pkg.NewAppError("Fail at bind search", err, http.StatusBadRequest)
	}
	users, err := eu.es.Search(&s)
	if err != nil {
		return err
	}
	if users == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, users)
}

func (eu EMPLuc) SearchWaiters(c echo.Context) error {
	s := model.Search{}
	if err := c.Bind(&s); err != nil {
		return pkg.NewAppError("Fail at bind search", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	users, err := eu.es.SearchWaiters(ur.EstablishmentID, &s)
	if err != nil {
		return err
	}
	if users == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, users)
}

func (eu EMPLuc) GetByID(c echo.Context) error {
	t, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	f, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	uj, err := eu.es.Get(f, uint(t))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, uj)
}

func (eu EMPLuc) Get(c echo.Context) error {
	id, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	uj, err := eu.es.Self(id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, uj)
}

func (eu EMPLuc) HireWaiter(c echo.Context) error {
	mail := c.Param("mail")
	r := model.UserRole{}
	if err := c.Bind(&r); err != nil {
		return pkg.NewAppError("Fail at bind role", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	if err = eu.es.HireWaiter(ur, mail, r.Salary); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Hire(c echo.Context) error {
	mail := c.Param("mail")
	r := model.UserRole{}
	if err := c.Bind(&r); err != nil {
		return pkg.NewAppError("Fail at bind role", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	if err = eu.es.Hire(ur, mail, &r); err != nil {
		return fmt.Errorf("fail at hire: %w", err)
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Fire(c echo.Context) error {
	t, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	var r struct {
		Reason string `json:"reason"`
	}
	if err := c.Bind(&r); err != nil {
		return pkg.NewAppError("Fail at bind reason", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	ur.Reason = r.Reason
	if err = eu.es.Fire(ur, uint(t)); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
