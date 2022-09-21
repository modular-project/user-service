package info

import (
	"context"
	"fmt"
	"net/http"
	"users-service/pkg"

	est "github.com/modular-project/protobuffers/information/establishment"
	"google.golang.org/grpc"
)

type establishmentService struct {
	esc est.EstablishmentServiceClient
}

func NewESTBService(conn *grpc.ClientConn) establishmentService {
	return establishmentService{esc: est.NewEstablishmentServiceClient(conn)}
}

func (e establishmentService) GetByAddress(ctx context.Context, aID string) (uint64, uint32, error) {
	r, err := e.esc.GetByAddress(ctx, &est.RequestGetByAddress{Id: aID})
	if err != nil {
		return 0, 0, fmt.Errorf("esc.GetByAddress: %w", err)
	}
	return r.Id, r.Quantity, nil
}

func (e establishmentService) Create(ctx context.Context, data *est.Establishment, qua uint32) (uint64, error) {
	res, err := e.esc.Create(ctx, &est.RequestCreate{Establishment: data, TableQuantity: qua})
	if err != nil {
		return 0, pkg.NewAppError("failed to create establishment", err, http.StatusInternalServerError)
	}
	return res.Id, nil
}

func (e establishmentService) GetByID(ctx context.Context, id uint64) (est.Establishment, error) {
	res, err := e.esc.Get(ctx, &est.RequestById{Id: id})
	if err != nil {
		return est.Establishment{}, pkg.NewAppError("failed to get establishment", err, http.StatusInternalServerError)
	}
	return *res, err
}

func (e establishmentService) GetInBatch(ctx context.Context, ids []uint64) ([]*est.Establishment, error) {
	res, err := e.esc.GetAll(ctx, &est.RequestGetAll{Ids: ids})
	if err != nil {
		return nil, pkg.NewAppError("failed to get establishment", err, http.StatusInternalServerError)
	}
	return res.Establishments, err
}

func (e establishmentService) Update(ctx context.Context, data *est.Establishment) error {
	_, err := e.esc.Update(ctx, &est.RequestUpdate{Establishment: data})
	if err != nil {
		return pkg.NewAppError("failed to update establishment", err, http.StatusInternalServerError)
	}
	return nil
}

func (e establishmentService) Delete(ctx context.Context, id uint64) error {
	_, err := e.esc.Delete(ctx, &est.RequestById{Id: id})
	if err != nil {
		return pkg.NewAppError("failed to delete establishment", err, http.StatusInternalServerError)
	}
	return nil
}
