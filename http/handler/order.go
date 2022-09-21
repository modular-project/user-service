package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"users-service/pkg"

	"github.com/labstack/echo"
	pf "github.com/modular-project/protobuffers/order/order"
)

type order struct {
	ID              uint64
	EstablishmentID uint64             `json:"establishment_id,omitempty"`
	Status          int32              `json:"status,omitempty"`
	Total           string             `json:"total,omitempty"`
	OrderProducts   []*pf.OrderProduct `json:"products,omitempty"`
	// Local
	EmployeeID uint64 `json:"employee_id,omitempty"`
	TableID    uint64 `json:"table_id,omitempty"`
	// Deliverty
	AddressID string `json:"address_id,omitempty"`
	UserID    uint64 `json:"user_id,omitempty"`
	CreatedAt uint64 `json:"created_at,omitempty"`
}

func newOrder(proto []*pf.Order) []order {
	if proto == nil {
		return nil
	}
	strconv.FormatFloat(13.0, 'f', 2, 32)
	o := make([]order, len(proto))
	for i := range proto {
		o[i] = order{
			ID:              proto[i].Id,
			EstablishmentID: proto[i].EstablishmentId,
			Status:          int32(proto[i].Status),
			Total:           strconv.FormatFloat(float64(proto[i].Total), 'f', 2, 64),
			OrderProducts:   proto[i].OrderProducts,
			CreatedAt:       proto[i].CreateAt,
		}
		if t := proto[i].GetLocalOrder(); t != nil {
			o[i].EmployeeID = t.EmployeeId
			o[i].TableID = t.TableId
		} else if t := proto[i].GetRemoteOrder(); t != nil {
			o[i].UserID = t.UserId
			o[i].AddressID = t.AddressId
		}
	}
	return o
}

type OrderServicer interface {
	CreateLocalOrder(ctx context.Context, in *pf.Order) (*pf.CreateResponse, float32, error)
	CreateDeliveryOrder(ctx context.Context, in *pf.Order) (*pf.CreateResponse, float32, error)
	GetOrdersByUser(ctx context.Context, s *pf.SearchOrders) (*pf.OrdersResponse, error)
	GetOrdersByKitchen(ctx context.Context, kID, last uint64) (*pf.OrderProductsResponse, error)
	GetOrders(ctx context.Context, s *pf.SearchOrders) (*pf.OrdersResponse, error)
	GetOrdersByEstablishment(ctx context.Context, s *pf.SearchOrders) (*pf.OrdersResponse, error)
	GetOrderByWaiter(ctx context.Context, wID uint64) (*pf.OrdersResponse, error)
	GetOrderByWaiterPending(ctx context.Context, wID uint64) (*pf.OrdersResponse, error)
	AddProductsToOrder(ctx context.Context, in *pf.AddProductsToOrderRequest) (*pf.AddProductsToOrderResponse, float32, error)
	GetOrderByID(ctx context.Context, id uint64) ([]*pf.OrderProduct, error)
}

type OrderUC struct {
	os OrderServicer
}

func NewOrderUC(os OrderServicer) OrderUC {
	return OrderUC{os: os}
}

func (ouc OrderUC) CreateLocalOrder(c echo.Context) error {
	var bo pf.Order
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	if err := c.Bind(&bo); err != nil {
		return pkg.NewAppError("Fail at bind local order", err, http.StatusBadRequest)
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	bo.EstablishmentId = uint64(ur.EstablishmentID)
	bo.Type = &pf.Order_LocalOrder{LocalOrder: &pf.LocalOrder{EmployeeId: uint64(ur.UserID), TableId: id}}
	r, t, err := ouc.os.CreateLocalOrder(c.Request().Context(), &bo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"order_id":    r.OrderId,
		"product_ids": r.ProductIds,
		"total":       t,
	})
}

func (ouc OrderUC) CreateDeliveryOrder(c echo.Context) error {
	var bo pf.Order
	if err := c.Bind(&bo); err != nil {
		return pkg.NewAppError("Fail at bind delivery order", err, http.StatusBadRequest)
	}
	uID, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	bo.Type = &pf.Order_RemoteOrder{RemoteOrder: &pf.RemoteOrder{UserId: uint64(uID)}}
	r, t, err := ouc.os.CreateDeliveryOrder(c.Request().Context(), &bo)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"order_id":    r.OrderId,
		"product_ids": r.ProductIds,
		"total":       t,
	})
}

func (ouc OrderUC) GetOrdersByUser(c echo.Context) error {
	var s pf.SearchOrders
	if err := c.Bind(&s); err != nil {
		return pkg.NewAppError("Fail at bind search", err, http.StatusBadRequest)
	}
	id, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	s.Users = []uint64{uint64(id)}
	r, err := ouc.os.GetOrdersByUser(c.Request().Context(), &s)
	if err != nil {
		return err
	}
	if r.Orders == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, newOrder(r.Orders))
}

func (ouc OrderUC) GetOrdersByKitchen(c echo.Context) error {
	id, err := getKitchenEstablishmentFromContext(c)
	if err != nil {
		return err
	}
	var last uint64
	if l := c.QueryParam("last"); l != "" {
		last, err = strconv.ParseUint(l, 10, 64)
		if err != nil {
			return fmt.Errorf("parse last: %w", err)
		}
	}
	r, err := ouc.os.GetOrdersByKitchen(c.Request().Context(), uint64(id), last)
	if err != nil {
		return err
	}
	log.Println("Antes: ", id, r.OrderProducts)
	if r.OrderProducts == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, r.OrderProducts)
}

func (ouc OrderUC) GetOrders(c echo.Context) error {
	var s pf.SearchOrders
	if err := c.Bind(&s); err != nil {
		return pkg.NewAppError("Fail at bind search", err, http.StatusBadRequest)
	}
	r, err := ouc.os.GetOrders(c.Request().Context(), &s)
	if err != nil {
		return err
	}
	if r.Orders == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, newOrder(r.Orders))
}

func (ouc OrderUC) GetOrdersByEstablishment(c echo.Context) error {
	var s pf.SearchOrders
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	if err := c.Bind(&s); err != nil {
		return pkg.NewAppError("Fail at bind search", err, http.StatusBadRequest)
	}
	s.Establishments = []uint64{
		uint64(ur.EstablishmentID),
	}
	r, err := ouc.os.GetOrdersByEstablishment(c.Request().Context(), &s)
	if err != nil {
		return err
	}
	if r.Orders == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, newOrder(r.Orders))
}
func (ouc OrderUC) GetOrderByWaiter(c echo.Context) error {
	id, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	r, err := ouc.os.GetOrderByWaiter(c.Request().Context(), uint64(id))
	if err != nil {
		return err
	}
	if r.Orders == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, newOrder(r.Orders))
}

func (ouc OrderUC) GetOrderByWaiterPending(c echo.Context) error {
	id, err := getUserIDFromContext(c)
	if err != nil {
		return err
	}
	r, err := ouc.os.GetOrderByWaiterPending(c.Request().Context(), uint64(id))
	if err != nil {
		return err
	}
	if r.Orders == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, newOrder(r.Orders))
}

func (ouc OrderUC) AddProductsToOrder(c echo.Context) error {
	var ps []*pf.OrderProduct
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	if err := c.Bind(&ps); err != nil {
		return pkg.NewAppError("Fail at bind order products", err, http.StatusBadRequest)
	}
	r, total, err := ouc.os.AddProductsToOrder(c.Request().Context(), &pf.AddProductsToOrderRequest{Id: id, Products: ps})
	if err != nil {
		return err
	}
	if r.Ids == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, map[string]interface{}{"ids": r.Ids, "total": total})
}

func (ouc OrderUC) GetProductsByOrderID(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	r, err := ouc.os.GetOrderByID(c.Request().Context(), id)
	if err != nil {
		return err
	}
	if r == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, r)
}
