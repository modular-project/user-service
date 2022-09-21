package controller

import (
	"net/http"
	"users-service/model"
	"users-service/pkg"
)

type PermissionStorage interface {
	UserRole(uint) (model.UserRole, error)
	Find(string) (model.User, error)
	IsVerified(uint) (bool, error)
	Kitchen(uint) (uint, error)
}

type permission struct {
	ps PermissionStorage
}

func NewPermission(ps PermissionStorage) permission {
	return permission{ps: ps}
}

func (p permission) Kitchen(kID uint) (uint, error) {
	aID, err := p.ps.Kitchen(kID)
	if err != nil {
		return 0, pkg.NewAppError("kitchen not found", err, http.StatusBadRequest)
	}
	return aID, nil
}

func (p permission) IsVerified(uID uint) (bool, error) {
	ok, err := p.ps.IsVerified(uID)
	if err != nil {
		return false, pkg.NewAppError("user not found", err, http.StatusBadRequest)
	}
	return ok, nil
}

func (p permission) UserRole(uID uint) (model.UserRole, error) {
	ur, err := p.ps.UserRole(uID)
	if err != nil {
		return model.UserRole{}, pkg.NewAppError("user role not found", err, http.StatusForbidden)
	}
	return ur, nil
}

func (p permission) Find(email string) (model.User, error) {
	u, err := p.ps.Find(email)
	if err != nil {
		return model.User{}, pkg.NewAppError("user not found", err, http.StatusBadRequest)
	}
	return u, nil
}

// func (p permission) Greater(uID uint, rID model.RoleID) error {
// 	u, err := p.Job(uID)
// 	if err != nil {
// 		return err
// 	}
// 	if !u.RoleID.IsGreater(rID) {
// 		return ErrDontHavePermission
// 	}
// 	return nil
// }

// // Equal Returns an error if the user does not have the role assigned, otherwise it returns the establishment ID with a null error
// func (p permission) Equal(uID uint, rID model.RoleID) (uint, error) {
// 	u, err := p.Job(uID)
// 	if err != nil {
// 		return 0, err
// 	}
// 	if u.RoleID != rID {
// 		return 0, ErrDontHavePermission
// 	}
// 	return u.EstablishmentID, nil
// }

// func (p permission) CanUpdate(from, target uint) error {
// 	f, err := p.Job(from)
// 	if err != nil {
// 		return err
// 	}
// 	if !f.RoleID.IsGreater(model.WAITER) {
// 		return ErrDontHavePermission
// 	}
// 	t, err := p.Job(target)
// 	if err != nil {
// 		return err
// 	}
// 	if t.RoleID == model.USER {
// 		return pkg.NewAppError("target is not a employee", nil, http.StatusBadRequest)
// 	}
// 	if !f.RoleID.IsGreater(t.RoleID) {
// 		return ErrDontHavePermission
// 	}
// 	if f.EstablishmentID == 0 || f.EstablishmentID == t.EstablishmentID {
// 		return nil
// 	}
// 	return ErrDontHavePermission
// }

// // CanHire return an error if contractor cannot hire the user or if user cannot be hired;
// // if error is nil it returns user id in *r
// // if contractor's establishment ID is non-zero, then set it in *r
// func (p permission) CanHire(con uint, email string, r *model.UserRole) error {
// 	if r.Salary <= 0 {
// 		return pkg.NewAppError("invalid salary", nil, http.StatusBadRequest)
// 	}
// 	// Verify if user can be an employee
// 	u, err := p.j.Find(email)
// 	if err != nil {
// 		return pkg.NewAppError("email not found", err, http.StatusBadRequest)
// 	}
// 	if u.IsActive {
// 		return pkg.NewAppError("target is already an employee", nil, http.StatusBadRequest)
// 	}
// 	if !u.IsVerified {
// 		return pkg.NewAppError("target is not verified", nil, http.StatusBadRequest)
// 	}
// 	r.UserID = u.ID
// 	// Verify if contractor can hire it
// 	conR, err := p.Job(con)
// 	if err != nil {
// 		return err
// 	}
// 	if !conR.RoleID.IsGreater(r.RoleID) {
// 		return ErrDontHavePermission
// 	}
// 	// Verify if the jobs need an establishment or not
// 	if conR.EstablishmentID != 0 {
// 		r.EstablishmentID = conR.EstablishmentID
// 	}
// 	need := p.needEstablishment(r.RoleID)
// 	if need && r.EstablishmentID == 0 {
// 		return pkg.NewAppError("role needs to have an establishment", nil, http.StatusForbidden)
// 	}
// 	if !need && r.EstablishmentID != 0 {
// 		return pkg.NewAppError("role need not have an establishment", nil, http.StatusForbidden)
// 	}
// 	return nil
// }
