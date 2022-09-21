package info

import (
	"context"
	"fmt"

	pf "github.com/modular-project/protobuffers/address/address"
	"google.golang.org/grpc"
)

// as addressService address.AddressServiceClient
type addressService struct {
	asc pf.AddressServiceClient
}

func NewAddressService(conn *grpc.ClientConn) addressService {
	return addressService{asc: pf.NewAddressServiceClient(conn)}
}

func (as addressService) CreateDelivery(ctx context.Context, d *pf.Delivery) (string, error) {
	r, err := as.asc.CreateDelivery(ctx, d)
	if err != nil {
		return "", fmt.Errorf("asc.CreateDelivery: %w", err)
	}
	return r.Id, nil
}

func (as addressService) GetAllByUser(ctx context.Context, u *pf.User) ([]*pf.Address, error) {
	r, err := as.asc.GetAllByUser(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("asc.GetAllByUser: %w", err)
	}
	return r.Address, nil
}

func (as addressService) DeleteByID(ctx context.Context, u *pf.User) (int64, error) {
	_, err := as.asc.DeleteByID(ctx, u)
	if err != nil {
		return 0, fmt.Errorf("asc.DeleteByID: %w", err)
	}
	return 1, nil
}

func (as addressService) GetAddByID(ctx context.Context, u *pf.ID) (*pf.Address, error) {
	r, err := as.asc.GetAddByID(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("asc.GetAddByID: %w", err)
	}
	return r, nil
}

func (as addressService) GetByID(ctx context.Context, u *pf.User) (*pf.Address, error) {
	r, err := as.asc.GetByID(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("asc.GetByID: %w", err)
	}
	return r, nil
}

func (as addressService) CreateEstablishment(ctx context.Context, a *pf.Address) (string, error) {
	r, err := as.asc.CreateEstablishment(ctx, a)
	if err != nil {
		return "", fmt.Errorf("asc.CreateEstablishment: %w", err)
	}
	return r.Id, nil
}

func (as addressService) DeleteEstablishment(ctx context.Context, aID string) (int64, error) {
	_, err := as.asc.DeleteEstablishment(ctx, &pf.ID{Id: aID})
	if err != nil {
		return 0, fmt.Errorf("asc.DeleteEstablishment: %w", err)
	}
	return 1, nil
}

func (as addressService) Search(ctx context.Context, sa *pf.SearchAddress) (*pf.ResponseAll, error) {
	r, err := as.asc.Search(ctx, sa)
	if err != nil {
		return nil, fmt.Errorf("asc.Search: %w", err)
	}
	return r, nil
}

func (as addressService) Nearest(ctx context.Context, uID uint64, aID string) (string, error) {
	r, err := as.asc.Nearest(ctx, &pf.User{Id: uID, AddressId: aID})
	if err != nil {
		return "", fmt.Errorf("asc.Nearest: %w", err)
	}
	return r.Id, nil
}
