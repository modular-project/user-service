package info

import (
	"context"
	"fmt"

	pf "github.com/modular-project/protobuffers/information/order"
	"google.golang.org/grpc"
)

type infoOrderService struct {
	clt pf.ValidateOrderClient
}

func NewInfoOrderService(conn *grpc.ClientConn) infoOrderService {
	return infoOrderService{clt: pf.NewValidateOrderClient(conn)}
}
func (ios infoOrderService) ValidateOrder(c context.Context, r *pf.ValidateOrderRequest) (float32, error) {
	res, err := ios.clt.ValidateOrder(c, r)
	if err != nil {
		return 0, fmt.Errorf("clt.ValidateOrder: %w", err)
	}
	return res.Total, nil
}
func (ios infoOrderService) ValidateProducts(c context.Context, r *pf.ValidateProductsRequest) (float32, error) {
	res, err := ios.clt.ValidateProducts(c, r)
	if err != nil {
		return 0, fmt.Errorf("clt.ValidateProducts: %w", err)
	}
	return res.Total, nil
}
