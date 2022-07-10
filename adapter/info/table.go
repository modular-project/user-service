package info

import (
	"context"
	"fmt"
	"net/http"
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
		return 0, fmt.Errorf("ts.pe.Job: %w", err)
	}
	// Check role and switch establishment ID
	if u.RoleID == model.MANAGER {
		eID = uint64(u.EstablishmentID)
	} else if !u.RoleID.IsGreater(model.MANAGER) {
		return 0, pkg.NewAppError("you don't have permission", nil, http.StatusForbidden)
	}
	if eID == 0 {
		return 0, pkg.NewAppError("establishment not found", nil, http.StatusBadRequest)
	}
	return eID, nil
}

// Delete qua tables in establishment with id = eID
// if user is a manager then eID is get from him
// returns the number of tables removed and an error if there are
func (ts tableService) Delete(ctx context.Context, uID uint, eID uint64, qua uint32) (uint32, error) {
	eID, err := ts.inEstablishment(uID, eID)
	if err != nil {
		return 0, fmt.Errorf("in establishment: %w", err)
	}
	r, err := ts.tc.RemoveFromEstablishment(ctx, &table.RequestDelete{EstablishmenId: eID, Quantity: qua})
	if err != nil {
		return 0, pkg.NewAppError("failed to remove tables", nil, http.StatusInternalServerError)
	}
	return r.Deleted, nil
}

// Deprecated, Use CreateInBatch
func (ts tableService) Create(ctx context.Context, uID uint, eID uint64) (uint64, error) {
	eID, err := ts.inEstablishment(uID, eID)
	if err != nil {
		return 0, fmt.Errorf("in establishment: %w", err)
	}
	// Call grpc service
	r, err := ts.tc.AddTable(ctx, &table.RequestById{Id: eID})
	if err != nil {
		return 0, pkg.NewAppError("failed to add table", err, http.StatusInternalServerError)
	}
	if r.Ids == nil {
		return 0, nil
	}
	return r.Ids[0], err
}

func (ts tableService) CreateInBatch(ctx context.Context, uID uint, eID uint64, qua uint32) ([]uint64, error) {
	eID, err := ts.inEstablishment(uID, eID)
	if err != nil {
		return nil, fmt.Errorf("in establishment: %w", err)
	}
	r, err := ts.tc.AddTables(ctx, &table.RequestAdd{Id: eID, Quantity: qua})
	if err != nil {
		return nil, pkg.NewAppError("failed to add tables", err, http.StatusInternalServerError)
	}
	return r.Ids, nil
}

func (ts tableService) GetFromEstablishment(ctx context.Context, eID uint64) ([]*table.Table, error) {
	r, err := ts.tc.GetFromEstablishment(ctx, &table.RequestById{Id: eID})
	if err != nil {
		return nil, pkg.NewAppError("failed to get tables", err, http.StatusInternalServerError)
	}
	return r.Tables, nil
}

// Not in use, a manager cannot change the status because they canceled the order in process
// Table status change automatically when an order is created or finished
func (ts tableService) ChangeStatus(ctx context.Context, t *table.Table) error {
	if t.Id == 0 {
		return pkg.NewAppError("table id not set", nil, http.StatusBadRequest)
	}
	u, err := ts.pe.Job(uint(t.UserId))
	if err != nil {
		return fmt.Errorf("ts.pe.Job: %w", err)
	}
	if u.EstablishmentID == 0 {
		return pkg.NewAppError("you have to be in an establishment", nil, http.StatusBadRequest)
	}
	t.EstablishmenId = uint64(u.EstablishmentID)
	_, err = ts.tc.ChangeStatus(ctx, t)
	if err != nil {
		return pkg.NewAppError("failed to change table status", err, http.StatusInternalServerError)
	}
	return nil
}
