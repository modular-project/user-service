package handler

import (
	"context"
	"net/http"
	"strconv"
	"users-service/pkg"

	"github.com/labstack/echo"
	"github.com/modular-project/protobuffers/information/product"
)

type ProductServicer interface {
	Create(context.Context, *product.Product) (uint64, error)
	Get(context.Context, uint64) (product.Product, error)
	GetAll(context.Context) ([]*product.Product, error)
	GetInBatch(context.Context, []uint64) ([]*product.Product, error)
	Delete(context.Context, uint64) error
	Update(context.Context, uint64, *product.Product) (uint64, error)
}

type ProductUC struct {
	ps ProductServicer
}

func NewProductUC(ps ProductServicer) ProductUC {
	return ProductUC{ps}
}

func (puc ProductUC) Create(c echo.Context) error {
	p := product.Product{}
	if err := c.Bind(&p); err != nil {
		return pkg.NewAppError("Fail at bind product", err, http.StatusBadRequest)
	}
	id, err := puc.ps.Create(c.Request().Context(), &p)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, responseID(id))
}

func (puc ProductUC) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	p, err := puc.ps.Get(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, p)
}

func (pub ProductUC) GetAll(c echo.Context) error {
	ps, err := pub.ps.GetAll(c.Request().Context())
	if err != nil {
		return err
	}
	if ps == nil {
		c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, ps)
}

func (pub ProductUC) GetInBatch(c echo.Context) error {
	type IDs struct {
		IDs []uint64 `json:"ids"`
	}
	var ids IDs
	if err := c.Bind(&ids); err != nil {
		return pkg.NewAppError("Fail at bind ids", err, http.StatusBadRequest)
	}
	ps, err := pub.ps.GetInBatch(c.Request().Context(), ids.IDs)
	if err != nil {
		return err
	}
	if ps == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, ps)
}

func (pub ProductUC) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	err = pub.ps.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (pub ProductUC) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	p := product.Product{}
	if err := c.Bind(&p); err != nil {
		return pkg.NewAppError("Fail at bind product", err, http.StatusBadRequest)
	}
	id, err = pub.ps.Update(c.Request().Context(), id, &p)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responseID(id))
}
