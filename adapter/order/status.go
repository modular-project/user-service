package order

import (
	"context"
	"fmt"

	pf "github.com/modular-project/protobuffers/order/order"
	"google.golang.org/grpc"
)

type orderStatusService struct {
	clt  pf.OrderStatusServiceClient
	near Nearester
	add  Addresser
}

type Nearester interface {
	Nearest(ctx context.Context, uID uint64, aID string) (string, error)
}

type Addresser interface {
	GetByAddress(context.Context, string) (uint64, uint32, error)
}

func NewOrderStatusService(conn *grpc.ClientConn, n Nearester, a Addresser) orderStatusService {
	return orderStatusService{clt: pf.NewOrderStatusServiceClient(conn), near: n, add: a}
}

func (oss orderStatusService) PayLocal(ctx context.Context, in *pf.PayLocalRequest) (*pf.PayLocalResponse, error) {
	r, err := oss.clt.PayLocal(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("clt.PayLocal: %w", err)
	}
	return r, nil
}
func (oss orderStatusService) PayDelivery(ctx context.Context, in *pf.PayDeliveryRequest) (*pf.PayDeliveryResponse, error) {
	eaID, err := oss.near.Nearest(ctx, in.UserId, in.Address)
	if err != nil {
		return nil, fmt.Errorf("near.Nearest: %w", err)
	}
	eID, _, err := oss.add.GetByAddress(ctx, eaID)
	if err != nil {
		return nil, fmt.Errorf("add.GetByAddress: %w", err)
	}
	in.EstablishmentId = eID
	r, err := oss.clt.PayDelivery(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("clt.PayDelivery: %w", err)
	}
	return r, nil
}
func (oss orderStatusService) CompleteProduct(ctx context.Context, in *pf.CompleteProductRequest) (*pf.CompleteProductResponse, error) {
	r, err := oss.clt.CompleteProduct(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("clt.CompleteProduct: %w", err)
	}
	return r, nil
}
func (oss orderStatusService) CapturePayment(ctx context.Context, in *pf.CapturePaymentRequest) (*pf.CapturePaymentResponse, error) {
	r, err := oss.clt.CapturePayment(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("clt.CapturePayment: %w", err)
	}
	return r, nil
}

func (oss orderStatusService) DeliverProduct(ctx context.Context, in *pf.DeliverProductRequest) (*pf.DeliverProductResponse, error) {
	r, err := oss.clt.DeliverProducts(ctx, in)
	if err != nil {
		return &pf.DeliverProductResponse{}, fmt.Errorf("clt.DeliverProducts: %w", err)
	}
	return r, nil
}
