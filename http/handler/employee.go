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
	Get(uint) (model.UserJobs, error)
	Search(*model.SearchEMPL) ([]model.User, error)
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
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at bind: %s", err)))
	}
	tID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get param id: %s", err)))
	}
	fID, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get id from context: %s", err)))
	}
	err = eu.es.Update(fID, uint(tID), u)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Search(c echo.Context) error {
	s := model.SearchEMPL{}
	if err := c.Bind(s); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at bind: %s", err)))
	}
	users, err := eu.es.Search(&s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at search: %s", err)))
	}
	if len(users) == 0 {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, users)
}

func (eu EMPLuc) SearchWaiters(c echo.Context) error {
	s := model.Search{}
	if err := c.Bind(s); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at bind: %s", err)))
	}
	eID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get param id: %s", err)))
	}
	users, err := eu.es.SearchWaiters(uint(eID), &s)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at search: %s", err)))
	}
	if len(users) == 0 {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, users)
}

func (eu EMPLuc) GetByID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get id from param: %s", err)))
	}
	uj, err := eu.es.Get(uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get user: %s", err)))
	}
	return c.JSON(http.StatusOK, uj)
}

func (eu EMPLuc) Get(c echo.Context) error {
	id, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get id from context: %s", err)))
	}
	uj, err := eu.es.Get(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get user: %s", err)))
	}
	return c.JSON(http.StatusOK, uj)
}

func (eu EMPLuc) HireWaiter(c echo.Context) error {
	mail := c.Param("mail")
	r := model.UserRole{}
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at bind: %s", err)))
	}
	id, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get id from context: %s", err)))
	}
	if err = eu.es.HireWaiter(id, mail, r.Salary); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at hire: %s", err)))
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Hire(c echo.Context) error {
	mail := c.Param("mail")
	r := model.UserRole{}
	if err := c.Bind(r); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at bind: %s", err)))
	}
	id, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get id from context: %s", err)))
	}
	if err = eu.es.Hire(id, mail, &r); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at hire: %s", err)))
	}
	return c.NoContent(http.StatusOK)
}

func (eu EMPLuc) Fire(c echo.Context) error {
	t, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get id from param: %s", err)))
	}
	f, err := getUserIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at get id from context: %s", err)))
	}
	if err = eu.es.Fire(f, uint(t)); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("fail at fire user: %s", err)))
	}
	return c.NoContent(http.StatusOK)
}
