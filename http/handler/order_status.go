package handler

import (
	"context"
	"net/http"
	"strconv"
	"users-service/pkg"

	"github.com/labstack/echo"
	pf "github.com/modular-project/protobuffers/order/order"
)

type OrderStatusServicer interface {
	PayLocal(ctx context.Context, in *pf.PayLocalRequest) (*pf.PayLocalResponse, error)
	PayDelivery(ctx context.Context, in *pf.PayDeliveryRequest) (*pf.PayDeliveryResponse, error)
	CompleteProduct(ctx context.Context, in *pf.CompleteProductRequest) (*pf.CompleteProductResponse, error)
	CapturePayment(ctx context.Context, in *pf.CapturePaymentRequest) (*pf.CapturePaymentResponse, error)
	DeliverProduct(context.Context, *pf.DeliverProductRequest) (*pf.DeliverProductResponse, error)
}

type OrderStatusUC struct {
	oss OrderStatusServicer
}

func NewOrderStatusUC(oss OrderStatusServicer) OrderStatusUC {
	return OrderStatusUC{oss: oss}
}

func (suc OrderStatusUC) PayLocal(c echo.Context) error {
	var p pf.PayLocalRequest
	if err := c.Bind(&p); err != nil {
		return pkg.NewAppError("failed to bind pay local request", err, http.StatusBadRequest)
	}
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	p.EmployeeId = uint64(uID)
	_, err = suc.oss.PayLocal(c.Request().Context(), &p)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (suc OrderStatusUC) PayDelivery(c echo.Context) error {
	var p pf.PayDeliveryRequest
	if err := c.Bind(&p); err != nil {
		return pkg.NewAppError("failed to bind pay delivery request", err, http.StatusBadRequest)
	}
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	p.UserId = uint64(uID)
	r, err := suc.oss.PayDelivery(c.Request().Context(), &p)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"id": r.Id})
}

func (suc OrderStatusUC) CompleteProduct(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	_, err = suc.oss.CompleteProduct(c.Request().Context(), &pf.CompleteProductRequest{Id: id})
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (suc OrderStatusUC) CapturePayment(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return pkg.NewAppError("Fail at get path param id", nil, http.StatusBadRequest)
	}
	r, err := suc.oss.CapturePayment(c.Request().Context(), &pf.CapturePaymentRequest{Id: id})
	if err != nil {
		return err
	}
	if r.Status == "" {
		return c.NoContent(http.StatusOK)
	}
	return c.String(http.StatusOK, r.Status)
}

func (suc OrderStatusUC) DeliverProducts(c echo.Context) error {
	var ids []uint64
	if err := c.Bind(&ids); err != nil {
		return pkg.NewAppError("failed to bind ids", err, http.StatusBadRequest)
	}
	if _, err := suc.oss.DeliverProduct(c.Request().Context(), &pf.DeliverProductRequest{Id: ids}); err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}
