package info

import (
	"context"
	"fmt"

	"github.com/modular-project/protobuffers/information/product"
	"google.golang.org/grpc"
)

type productService struct {
	pc product.ProductServiceClient
}

func NewProductService(conn *grpc.ClientConn) productService {
	return productService{pc: product.NewProductServiceClient(conn)}
}

func (ps productService) Create(ctx context.Context, p *product.Product) (uint64, error) {
	r, err := ps.pc.Create(ctx, p)
	if err != nil {
		return 0, err
	}
	return r.Id, nil
}

func (ps productService) Get(ctx context.Context, id uint64) (product.Product, error) {
	r, err := ps.pc.Get(ctx, &product.RequestById{Id: id})
	if err != nil {
		return product.Product{}, err
	}
	return *r, err
}

func (ps productService) GetAll(ctx context.Context) ([]*product.Product, error) {
	r, err := ps.pc.GetAll(ctx, &product.RequestGetAll{})
	if err != nil {
		return nil, err
	}
	return r.Products, err
}

func (ps productService) GetInBatch(ctx context.Context, IDs []uint64) ([]*product.Product, error) {
	r, err := ps.pc.GetInBatch(ctx, &product.RequestGetInBatch{Ids: IDs})
	if err != nil {
		return nil, fmt.Errorf("get in batch: %w", err)
	}
	return r.Products, err
}

func (ps productService) Delete(ctx context.Context, id uint64) error {
	_, err := ps.pc.Delete(ctx, &product.RequestById{Id: id})
	return err
}

func (ps productService) Update(ctx context.Context, id uint64, p *product.Product) (uint64, error) {
	r, err := ps.pc.Update(ctx, &product.RequestUpdate{Id: id, Product: p})
	if err != nil {
		return 0, err
	}
	return r.Id, nil
}
