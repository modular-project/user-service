package handler

import (
	"fmt"
	"net/http"
	"strconv"
	"users-service/model"

	"github.com/labstack/echo"
)

type EMPLService interface {
	Update(from, target uint, u *model.User) error
	Self(uint) (model.UserJobs, error)
	Get(from, target uint) (model.UserJobs, error)
	Search(uint, *model.SearchEMPL) ([]model.User, error)
	SearchWaiters(uint, *model.Search) ([]model.User, error)
	Hire(uint, string, *model.UserRole) error
	HireWaiter(uint, string, float64) error
	Fire(from uint, target uint) error
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
		return fmt.Errorf("%w, %v", ErrBindData, err)
	}
	tID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetParamFromPath, err)
	}
	fID, err := getUserIDFromContext(c)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetIDFromContext, err)
	}
	err = eu.es.Update(fID, uint(tID), u)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Search(c echo.Context) error {
	s := model.SearchEMPL{}
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetIDFromContext, err)
	}
	if err := c.Bind(&s); err != nil {
		return fmt.Errorf("%w searchEMPL, %v", ErrBindData, err)
	}
	users, err := eu.es.Search(uID, &s)
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
		return fmt.Errorf("%w search, %v", ErrBindData, err)
	}
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetIDFromContext, err)
	}
	users, err := eu.es.SearchWaiters(uID, &s)
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
		return fmt.Errorf("%w, %v", ErrGetParamFromPath, err)
	}
	f, err := getUserIDFromContext(c)
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
		return fmt.Errorf("%w, %v", ErrGetIDFromContext, err)
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
		return fmt.Errorf("%w, %v", ErrBindData, err)
	}
	id, err := getUserIDFromContext(c)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetIDFromContext, err)
	}
	if err = eu.es.HireWaiter(id, mail, r.Salary); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Hire(c echo.Context) error {
	mail := c.Param("mail")
	r := model.UserRole{}
	if err := c.Bind(&r); err != nil {
		return fmt.Errorf("%w, %v", ErrBindData, err)
	}
	id, err := getUserIDFromContext(c)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetIDFromContext, err)
	}
	if err = eu.es.Hire(id, mail, &r); err != nil {
		return fmt.Errorf("fail at hire: %w", err)
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Fire(c echo.Context) error {
	t, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetParamFromPath, err)
	}
	f, err := getUserIDFromContext(c)
	if err != nil {
		return fmt.Errorf("%w, %v", ErrGetIDFromContext, err)
	}
	if err = eu.es.Fire(f, uint(t)); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
