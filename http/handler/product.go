package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

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
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	id, err := puc.ps.Create(context.TODO(), &p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	return c.JSON(http.StatusCreated, responseID(id))
}

func (puc ProductUC) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	p, err := puc.ps.Get(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, p)
}

func (pub ProductUC) GetAll(c echo.Context) error {
	ps, err := pub.ps.GetAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, ps)
}

func (pub ProductUC) GetInBatch(c echo.Context) error {
	type IDs struct {
		IDs []uint64 `json:"ids"`
	}
	var ids IDs
	if err := c.Bind(&ids); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("bind: %s", err)))
	}
	ps, err := pub.ps.GetInBatch(context.Background(), ids.IDs)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, ps)
}

func (pub ProductUC) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	err = pub.ps.Delete(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	return c.NoContent(http.StatusOK)
}

func (pub ProductUC) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	p := product.Product{}
	if err := c.Bind(&p); err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	id, err = pub.ps.Update(context.Background(), id, &p)
	if err != nil {
		return c.JSON(http.StatusBadRequest, createResponse(err.Error()))
	}
	return c.JSON(http.StatusOK, responseID(id))
}
