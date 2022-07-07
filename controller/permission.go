package controller

import (
	"fmt"
	"users-service/model"
)

type JobStorager interface {
	Job(uint) (model.UserRole, error)
	Find(string) (model.User, error)
	IsVerified(uint) (bool, error)
}

type permission struct {
	j JobStorager
}

func NewPermission(j JobStorager) permission {
	return permission{j: j}
}

func (p permission) IsVerified(uID uint) (bool, error) {
	return p.j.IsVerified(uID)
}

func (p permission) needEstablishment(role model.RoleID) bool {
	switch role {
	case model.ADMIN, model.OWNER:
		return false
	}
	return true
}

func (p permission) Greater(uID uint, rID model.RoleID) error {
	u, err := p.j.Job(uID)
	if err != nil {
		return err
	}
	if !u.RoleID.IsGreater(rID) {
		return ErrUnauthorizedUser
	}
	return nil
}

// Equal Returns an error if the user does not have the role assigned, otherwise it returns the establishment ID with a null error
func (p permission) Equal(uID uint, rID model.RoleID) (uint, error) {
	u, err := p.j.Job(uID)
	if err != nil {
		return 0, err
	}
	if u.RoleID != rID {
		return 0, ErrUnauthorizedUser
	}
	return u.EstablishmentID, nil
}

func (p permission) CanUpdate(from, target uint) error {
	f, err := p.j.Job(from)
	if err != nil {
		return err
	}
	if !f.RoleID.IsGreater(model.WAITER) {
		return ErrUnauthorizedUser
	}
	t, err := p.j.Job(target)
	if err != nil {
		return err
	}
	if t.RoleID == model.USER {
		return ErrIsNotAnEmployee
	}
	if !f.RoleID.IsGreater(t.RoleID) {
		return ErrUnauthorizedUser
	}
	if f.EstablishmentID == 0 || f.EstablishmentID == t.EstablishmentID {
		return nil
	}
	return ErrUnauthorizedUser
}

// CanHire return an error if contractor cannot hire the user or if user cannot be hired;
// if error is nil it returns user id in *r
// if contractor's establishment ID is non-zero, then set it in *r
func (p permission) CanHire(con uint, email string, r *model.UserRole) error {
	if r.Salary <= 0 {
		return ErrInvalidSalary
	}
	// Verify if user can be an employee
	u, err := p.j.Find(email)
	if err != nil {
		return fmt.Errorf("find mail %s: %w", email, err)
	}
	if u.IsActive {
		return ErrAlreadyEmployee
	}
	if !u.IsVerified {
		return ErrUserIsNotVerified
	}
	r.UserID = u.ID
	// Verify if contractor can hire it
	conR, err := p.j.Job(con)
	if err != nil {
		return err
	}
	if !conR.RoleID.IsGreater(r.RoleID) {
		return ErrUnauthorizedUser
	}
	// Verify if the jobs need an establishment or not
	if conR.EstablishmentID != 0 {
		r.EstablishmentID = conR.EstablishmentID
	}
	need := p.needEstablishment(r.RoleID)
	if need && r.EstablishmentID == 0 {
		return ErrEstablishNecesary
	}
	if !need && r.EstablishmentID != 0 {
		return ErrCannotBeAssigned
	}
	return nil
}
