package info

import (
	"context"
	"fmt"
	"users-service/model"
	"users-service/pkg"

	"github.com/modular-project/protobuffers/information/table"
	"google.golang.org/grpc"
)

type Permissioner interface {
	Job(uint) (model.UserRole, error)
}

type tableService struct {
	tc table.TableServiceClient
	pe Permissioner
}

func NewTableService(conn *grpc.ClientConn, pe Permissioner) tableService {
	return tableService{tc: table.NewTableServiceClient(conn), pe: pe}
}

// inEstablishment Return an error if the user cannot modify the setting
// If the user is a manager, then it returns their establishment ID
// If the eID is zero, it returns an error.
func (ts tableService) inEstablishment(uID uint, eID uint64) (uint64, error) {
	u, err := ts.pe.Job(uID)
	if err != nil {
		return 0, err
	}
	// Check role and switch establishment ID
	if u.RoleID == model.MANAGER {
		eID = uint64(u.EstablishmentID)
	} else if !u.RoleID.IsGreater(model.MANAGER) {
		return 0, pkg.BadErr("you don't have the necessary role")
	}
	if eID == 0 {
		return 0, pkg.NotFoundErr("establishment not found")
	}
	return eID, nil
}

// Delete qua tables in establishment with id = eID
// if user is a manager then eID is get from him
// returns the number of tables removed and an error if there are
func (ts tableService) Delete(ctx context.Context, uID uint, eID uint64, qua uint32) (uint32, error) {
	eID, err := ts.inEstablishment(uID, eID)
	if err != nil {
		return 0, err
	}
	r, err := ts.tc.RemoveFromEstablishment(ctx, &table.RequestDelete{EstablishmenId: eID, Quantity: qua})
	if err != nil {
		return 0, fmt.Errorf("remove from est: %w", err)
	}
	return r.Deleted, nil
}

// Deprecated, Use CreateInBatch
func (ts tableService) Create(ctx context.Context, uID uint, eID uint64) (uint64, error) {
	eID, err := ts.inEstablishment(uID, eID)
	if err != nil {
		return 0, err
	}
	// Call grpc service
	r, err := ts.tc.AddTable(ctx, &table.RequestById{Id: eID})
	if err != nil {
		return 0, fmt.Errorf("add table: %w", err)
	}
	if r.Ids == nil {
		return 0, nil
	}
	return r.Ids[0], err
}

func (ts tableService) CreateInBatch(ctx context.Context, uID uint, eID uint64, qua uint32) ([]uint64, error) {
	eID, err := ts.inEstablishment(uID, eID)
	if err != nil {
		return nil, err
	}
	r, err := ts.tc.AddTables(ctx, &table.RequestAdd{Id: eID, Quantity: qua})
	if err != nil {
		return nil, fmt.Errorf("add tables: %w", err)
	}
	if r.Ids == nil {
		return nil, nil
	}
	return r.Ids, nil
}

func (ts tableService) GetFromEstablishment(ctx context.Context, eID uint64) ([]*table.Table, error) {
	r, err := ts.tc.GetFromEstablishment(ctx, &table.RequestById{Id: eID})
	if err != nil {
		return nil, fmt.Errorf("get tables: %w", err)
	}
	if r.Tables == nil {
		return nil, nil
	}
	return r.Tables, nil
}

// Not in use, a manager cannot change the status because they canceled the order in process
// Table status change automatically when an order is created or finished
func (ts tableService) ChangeStatus(ctx context.Context, t *table.Table) error {
	if t.Id == 0 {
		return pkg.BadErr("table id not set")
	}
	u, err := ts.pe.Job(uint(t.UserId))
	if err != nil {
		return err
	}
	if u.EstablishmentID == 0 {
		return pkg.ForbiddenErr("you have to be in an establishment")
	}
	t.EstablishmenId = uint64(u.EstablishmentID)
	_, err = ts.tc.ChangeStatus(ctx, t)
	if err != nil {
		return fmt.Errorf("change table status: %w", err)
	}
	return nil
}
