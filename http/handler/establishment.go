package handler

import (
	"context"
	"net/http"
	"strconv"
	"users-service/pkg"

	"github.com/labstack/echo"
	est "github.com/modular-project/protobuffers/information/establishment"
)

type ESTDService interface {
	Create(context.Context, *est.Establishment, uint32) (uint64, error)
	GetByID(context.Context, uint64) (est.Establishment, error)
	GetInBatch(context.Context, []uint64) ([]*est.Establishment, error)
	Update(context.Context, *est.Establishment) error
	Delete(context.Context, uint64) error
}

type ESTDuc struct {
	ess ESTDService
}

func NewESTDuc(ess ESTDService) ESTDuc {
	return ESTDuc{ess: ess}
}

func (euc ESTDuc) Create(c echo.Context) error {
	q, err := strconv.ParseUint(c.QueryParam("q"), 10, 0)
	if err != nil {
		q = 1
	}
	e := est.Establishment{}
	if err := c.Bind(&e); err != nil {
		return pkg.NewAppError("fail at bind establishment", err, http.StatusBadRequest)
	}
	id, err := euc.ess.Create(c.Request().Context(), &e, uint32(q))
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, responseID(id))
}

func (euc ESTDuc) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	es, err := euc.ess.GetByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, es)
}

func (euc ESTDuc) GetInBatch(c echo.Context) error {
	type IDs struct {
		IDs []uint64 `json:"ids"`
	}
	var ids IDs
	if err := c.Bind(&ids); err != nil {
		return pkg.NewAppError("Fail at bind ids", err, http.StatusBadRequest)
	}
	es, err := euc.ess.GetInBatch(c.Request().Context(), ids.IDs)
	if err != nil {
		return err
	}
	if es == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, es)
}

func (euc ESTDuc) Update(c echo.Context) error {
	var es est.Establishment
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	if err := c.Bind(&es); err != nil {
		return pkg.NewAppError("Fail at bind ids", err, http.StatusBadRequest)
	}
	es.Id = id
	if err := euc.ess.Update(c.Request().Context(), &es); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (euc ESTDuc) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	if err := euc.ess.Delete(c.Request().Context(), id); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
