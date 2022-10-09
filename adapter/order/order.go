package order

import (
	"context"
	"fmt"

	val "github.com/modular-project/protobuffers/information/order"
	pf "github.com/modular-project/protobuffers/order/order"
	"google.golang.org/grpc"
)

type validater interface {
	ValidateOrder(context.Context, *val.ValidateOrderRequest) (float32, error)
	ValidateProducts(context.Context, *val.ValidateProductsRequest) (float32, error)
}

type orderService struct {
	clt pf.OrderServiceClient
	v   validater
}

func NewOrderService(conn *grpc.ClientConn, v validater) orderService {
	return orderService{clt: pf.NewOrderServiceClient(conn), v: v}
}

// Return true if an establishment have pending orders
func (os orderService) HavePendingOrders(ctx context.Context, eID uint64) (bool, error) {
	r, err := os.clt.GetOrdersByEstablishment(ctx, &pf.OrdersRequest{Search: &pf.SearchOrders{Default: &pf.Default{Limit: 1}, Establishments: []uint64{eID}, Status: []pf.Status{2}}})
	if err != nil {
		return false, fmt.Errorf("clt.GetOrdersByEstablishment: %w", err)
	}
	return r.Orders != nil, nil
}

func (os orderService) CreateLocalOrder(ctx context.Context, in *pf.Order) (*pf.CreateResponse, float32, error) {
	t, err := os.v.ValidateOrder(ctx, &val.ValidateOrderRequest{Order: in})
	if err != nil {
		return nil, 0, fmt.Errorf("validate Order: %w", err)
	}
	in.Total = t
	r, err := os.clt.CreateLocalOrder(ctx, in)
	if err != nil {
		return nil, 0, fmt.Errorf("createLocalOrder: %w", err)
	}
	return r, t, nil
}

func (os orderService) CreateDeliveryOrder(ctx context.Context, in *pf.Order) (*pf.CreateResponse, float32, error) {
	t, err := os.v.ValidateProducts(ctx, &val.ValidateProductsRequest{OrderProducts: in.OrderProducts})
	if err != nil {
		return nil, 0, fmt.Errorf("validate Order: %w", err)
	}
	in.Total = t
	r, err := os.clt.CreateDeliveryOrder(ctx, in)
	if err != nil {
		return nil, 0, fmt.Errorf("createDeliveryOrder: %w", err)
	}
	return r, t, nil
}

func (os orderService) GetOrdersByUser(ctx context.Context, in *pf.SearchOrders) (*pf.OrdersResponse, error) {
	if in == nil {
		return nil, fmt.Errorf("nil search")
	}
	o, err := os.clt.GetOrdersByUser(ctx, &pf.OrdersByUserRequest{Search: in})
	if err != nil {
		return nil, fmt.Errorf("clt.GetOrdersByUser: %w", err)
	}
	return o, nil
}

func (os orderService) GetOrdersByKitchen(ctx context.Context, kID uint64, l uint64) (*pf.OrderProductsResponse, error) {
	r, err := os.clt.GetOrdersByKitchen(ctx, &pf.RequestKitchen{Id: kID, Last: l})
	if err != nil {
		return nil, fmt.Errorf("clt.GetOrdersByKitchen: %w", err)
	}
	return r, nil
}

func (os orderService) GetOrders(ctx context.Context, s *pf.SearchOrders) (*pf.OrdersResponse, error) {
	r, err := os.clt.GetOrders(ctx, &pf.OrdersRequest{Search: s})
	if err != nil {
		return nil, fmt.Errorf("clt.GetOrders: %w", err)
	}
	return r, nil
}

func (os orderService) GetOrdersByEstablishment(ctx context.Context, s *pf.SearchOrders) (*pf.OrdersResponse, error) {
	r, err := os.clt.GetOrdersByEstablishment(ctx, &pf.OrdersRequest{Search: s})
	if err != nil {
		return nil, fmt.Errorf("clt.GetOrdersByEstablishment: %w", err)
	}
	return r, nil
}

func (os orderService) GetOrderByWaiter(ctx context.Context, wID uint64) (*pf.OrdersResponse, error) {
	r, err := os.clt.GetOrderByWaiter(ctx, &pf.ID{Id: wID})
	if err != nil {
		return nil, fmt.Errorf("clt.GetOrdersByWaiter: %w", err)
	}
	return r, nil
}

func (os orderService) GetOrderByWaiterPending(ctx context.Context, wID uint64) (*pf.OrdersResponse, error) {
	r, err := os.clt.GetOrderPendingByWaiter(ctx, &pf.ID{Id: wID})
	if err != nil {
		return nil, fmt.Errorf("clt.GetOrdersByWaiter: %w", err)
	}
	return r, nil
}

func (os orderService) AddProductsToOrder(ctx context.Context, in *pf.AddProductsToOrderRequest) (*pf.AddProductsToOrderResponse, float32, error) {
	t, err := os.v.ValidateProducts(ctx, &val.ValidateProductsRequest{OrderProducts: in.Products})
	if err != nil {
		return nil, 0, fmt.Errorf("v.ValidateProducts: %w", err)
	}
	in.Total = t
	r, err := os.clt.AddProductsToOrder(ctx, in)
	if err != nil {
		return nil, 0, fmt.Errorf("clt.AddProductsToOrder: %w", err)
	}
	return r, t, nil
}

func (os orderService) GetOrderByID(ctx context.Context, id uint64) ([]*pf.OrderProduct, error) {
	r, err := os.clt.GetOrderByID(ctx, &pf.GetOrderByIDRequest{OrderId: id})
	if err != nil {
		return nil, fmt.Errorf("clt.GetOrderByID: %w", err)
	}
	return r.Order.OrderProducts, nil
}
