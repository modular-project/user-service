package handler

import (
	"context"
	"net/http"
	"users-service/pkg"

	"github.com/labstack/echo"
	apf "github.com/modular-project/protobuffers/address/address"
)

type AddDeliveryUC struct {
	ads AddDeliveryServicer
}

type AddDeliveryServicer interface {
	CreateDelivery(context.Context, *apf.Delivery) (string, error)
	GetAllByUser(context.Context, *apf.User) ([]*apf.Address, error)
	//Nearest(context.Context, *apf.User) (string, error) // TODO: ADD TO ORDERS
	DeleteByID(context.Context, *apf.User) (int64, error)
	GetByID(context.Context, *apf.User) (*apf.Address, error)
}

func NewAddDeliveryUC(ads AddDeliveryServicer) AddDeliveryUC {
	return AddDeliveryUC{ads: ads}
}

func (auc AddDeliveryUC) Create(c echo.Context) error {
	var d apf.Delivery
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}

	if err := c.Bind(&d.Address); err != nil {
		return pkg.NewAppError("fail at bind address", err, http.StatusBadRequest)
	}
	d.UserId = uint64(uID)
	aID, err := auc.ads.CreateDelivery(c.Request().Context(), &d)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"id": aID})
}

func (auc AddDeliveryUC) GetAllByUser(c echo.Context) error {
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}

	adds, err := auc.ads.GetAllByUser(c.Request().Context(), &apf.User{Id: uint64(uID)})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, adds)
}

func (auc AddDeliveryUC) Delete(c echo.Context) error {
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	aID := c.Param("id")
	if aID == "" {
		return pkg.NewAppError("Fail at get path param address id", err, http.StatusBadRequest)
	}
	if _, err := auc.ads.DeleteByID(c.Request().Context(), &apf.User{Id: uint64(uID), AddressId: aID}); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (auc AddDeliveryUC) GetByID(c echo.Context) error {
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	aID := c.Param("id")
	if aID == "" {
		return pkg.NewAppError("Fail at get path param address id", err, http.StatusBadRequest)
	}
	a, err := auc.ads.GetByID(c.Request().Context(), &apf.User{Id: uint64(uID), AddressId: aID})
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, a)
}
