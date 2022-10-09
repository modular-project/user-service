package handler

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"users-service/pkg"

	"github.com/labstack/echo"
	apf "github.com/modular-project/protobuffers/address/address"
	est "github.com/modular-project/protobuffers/information/establishment"
)

type ESTDService interface {
	Create(context.Context, *est.Establishment, uint32) (uint64, error)
	GetByID(context.Context, uint64) (est.Establishment, error)
	GetInBatch(context.Context, []uint64) ([]*est.Establishment, error)
	Update(context.Context, *est.Establishment) error
	Delete(context.Context, uint64) (string, error) // Return ID address from deleted
	GetByAddress(context.Context, string) (uint64, uint32, error)
}

type AddressServicer interface {
	CreateEstablishment(context.Context, *apf.Address) (string, error)
	DeleteEstablishment(context.Context, string) (int64, error)
	//Search(context.Context, *apf.SearchAddress) ([]*apf.Address, error)
	GetAddByID(context.Context, *apf.ID) (*apf.Address, error)
	Search(context.Context, *apf.SearchAddress) (*apf.ResponseAll, error)
}

type HavePendingOrderser interface {
	HavePendingOrders(context.Context, uint64) (bool, error)
}

type HaveEMPLser interface {
	HaveActiveEMPLs(uint) (bool, error)
}

type ESTDuc struct {
	ess ESTDService
	as  AddressServicer
	ho  HavePendingOrderser
	he  HaveEMPLser
}

func NewESTDuc(ess ESTDService, as AddressServicer, ho HavePendingOrderser, he HaveEMPLser) ESTDuc {
	return ESTDuc{ess: ess, as: as, ho: ho, he: he}
}

func (euc ESTDuc) Create(c echo.Context) error {
	q, err := strconv.ParseUint(c.QueryParam("q"), 10, 0)
	if err != nil {
		q = 1
	}
	add := apf.Address{}
	if err := c.Bind(&add); err != nil {
		return pkg.NewAppError("fail at bind establishment", err, http.StatusBadRequest)
	}
	aID, err := euc.as.CreateEstablishment(c.Request().Context(), &add)
	if err != nil {
		return err
	}
	log.Print(aID)
	e := est.Establishment{AddressId: aID}
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
	add, err := euc.as.GetAddByID(c.Request().Context(), &apf.ID{Id: es.AddressId})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"address": add, "quantity": es.Quantity})
}

func (euc ESTDuc) GetInBatch(c echo.Context) error {
	type IDs struct {
		IDs []uint64 `json:"ids"`
	}
	// TODO: ADD ADDRESS IN RESULT
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
	// TODO: UPDATE ADDRESS
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
	have, err := euc.ho.HavePendingOrders(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if have {
		return pkg.NewAppError("Establishment have pending orders", nil, http.StatusBadRequest)
	}
	he, err := euc.he.HaveActiveEMPLs(uint(id))
	if err != nil {
		return err
	}
	if he {
		return pkg.NewAppError("Establishment have employees", nil, http.StatusBadRequest)
	}
	aID, err := euc.ess.Delete(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if _, err := euc.as.DeleteEstablishment(c.Request().Context(), aID); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (euc ESTDuc) Search(c echo.Context) error {
	var s apf.SearchAddress
	if err := c.Bind(&s); err != nil {
		return pkg.NewAppError("Fail at bind search", err, http.StatusBadRequest)
	}
	es, err := euc.as.Search(c.Request().Context(), &s)
	if err != nil {
		return err
	}
	if es == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, es.Address)
}

func (euc ESTDuc) GetByAddress(c echo.Context) error {
	aID := c.Param("id")
	if aID == "" {
		return pkg.NewAppError("Fail at get path param id", nil, http.StatusBadRequest)
	}
	eID, q, err := euc.ess.GetByAddress(c.Request().Context(), aID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"id": eID, "quantity": q})
}
