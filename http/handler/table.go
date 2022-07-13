package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"users-service/pkg"

	"github.com/labstack/echo"
	"github.com/modular-project/protobuffers/information/table"
)

type TableService interface {
	Delete(ctx context.Context, eID uint64, qua uint32) (uint32, error)
	CreateInBatch(ctx context.Context, eID uint64, qua uint32) ([]uint64, error)
	GetFromEstablishment(context.Context, uint64) ([]*table.Table, error)
}

type TableUC struct {
	ts TableService
}

func NewTableUC(ts TableService) TableUC {
	return TableUC{ts}
}

func (tu TableUC) Create(c echo.Context) error {
	eID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	q, err := strconv.ParseUint(c.QueryParam("q"), 10, 0)
	if err != nil {
		q = 1
	}
	ids, err := tu.ts.CreateInBatch(c.Request().Context(), eID, uint32(q))
	if err != nil {
		return fmt.Errorf("fail at create: %w", err)
	}
	if ids == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, ids)
}

func (tu TableUC) CreateIn(c echo.Context) error {
	q, err := strconv.ParseUint(c.QueryParam("q"), 10, 0)
	if err != nil {
		q = 1
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}

	ids, err := tu.ts.CreateInBatch(c.Request().Context(), uint64(ur.EstablishmentID), uint32(q))
	if err != nil {
		return fmt.Errorf("fail at create: %w", err)
	}
	if ids == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, ids)
}

func (tu TableUC) Delete(c echo.Context) error {
	eID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	q, err := strconv.ParseUint(c.QueryParam("q"), 10, 0)
	if err != nil {
		q = 1
	}
	del, err := tu.ts.Delete(c.Request().Context(), eID, uint32(q))
	if err != nil {
		return fmt.Errorf("fail at delete: %w", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"deleted": del})
}

func (tu TableUC) DeleteIn(c echo.Context) error {
	q, err := strconv.ParseUint(c.QueryParam("q"), 10, 0)
	if err != nil {
		q = 1
	}
	ur, err := getUserRoleFromContext(c)
	if err != nil {
		return err
	}
	del, err := tu.ts.Delete(c.Request().Context(), uint64(ur.EstablishmentID), uint32(q))
	if err != nil {
		return fmt.Errorf("fail at delete: %w", err)
	}

	return c.JSON(http.StatusOK, echo.Map{"deleted": del})
}

func (tu TableUC) Get(c echo.Context) error {
	eID, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		return pkg.NewAppError("Fail at get path param id", err, http.StatusBadRequest)
	}
	ta, err := tu.ts.GetFromEstablishment(c.Request().Context(), eID)
	if err != nil {
		return fmt.Errorf("get from est: %w", err)
	}
	if ta == nil {
		return c.NoContent(http.StatusOK)
	}
	return c.JSON(http.StatusOK, ta)
}
