package controller

import (
	"net/http"
	"users-service/model"
	"users-service/pkg"
)

var (
	ErrDontHavePermission = pkg.NewAppError("you don't have permission", nil, http.StatusForbidden)
)

type EMPLStorager interface {
	Self(uint) (model.UserJobs, error)
	Get(from *model.UserRole, target uint) (model.UserJobs, error)
	SearchWaiters(uint, *model.Search) ([]model.User, error)
	Search(*model.SearchEMPL) ([]model.User, error)
	Hire(*model.UserRole) error
	Fire(uint) error
}

type Canner interface {
	CanUpdate(from, target uint) error
	CanHire(uint, string, *model.UserRole) error
	Greater(uID uint, rID model.RoleID) error
	Equal(uID uint, rID model.RoleID) (uint, error)
	Job(uint) (model.UserRole, error)
}

type Finder interface {
	Find(string) (model.User, error)
	UserRole(uint) (model.UserRole, error)
}

type Updater interface {
	Update(*model.User) error
}

type employeeService struct {
	est EMPLStorager
	up  Updater
	fi  Finder
}

// type EMPLService interface {
// 	Update(from model.UserRole, target uint, u *model.User) error
// 	Self(uint) (model.UserJobs, error)
// 	Get(from model.UserRole, target uint) (model.UserJobs, error)
// 	Search(*model.SearchEMPL) ([]model.User, error)
// 	SearchWaiters(uint, *model.Search) ([]model.User, error)
// 	Hire(model.UserRole, string, *model.UserRole) error
// 	HireWaiter(model.UserRole, string, float64) error
// 	Fire(from model.UserRole, target uint) error
// }

func NewEmployeeService(est EMPLStorager, up Updater, fi Finder) employeeService {
	return employeeService{est: est, up: up, fi: fi}
}

func (es employeeService) needEstablishment(role model.RoleID) bool {
	switch role {
	case model.ADMIN, model.OWNER:
		return false
	}
	return true
}

func (es employeeService) canUpdate(f model.UserRole, target uint) error {
	t, err := es.fi.UserRole(target)
	if err != nil {
		return err
	}
	if t.RoleID == model.USER {
		return pkg.NewAppError("target is not a employee", nil, http.StatusBadRequest)
	}
	if !f.RoleID.IsGreater(t.RoleID) {
		return ErrDontHavePermission
	}

	if f.EstablishmentID == 0 || f.EstablishmentID == t.EstablishmentID {
		return nil
	}
	return ErrDontHavePermission
}

// CanHire return an error if contractor cannot hire the user or if user cannot be hired;
// if error is nil it returns user id in *r
// if contractor's establishment ID is non-zero, then set it in *r
func (es employeeService) canHire(f, tr *model.UserRole, email string) error {
	if tr.Salary <= 0 {
		return pkg.NewAppError("salary must be a positive number", nil, http.StatusBadRequest)
	}
	// Verify if user can be an employee
	u, err := es.fi.Find(email)
	if err != nil {
		return pkg.NewAppError("email not found", err, http.StatusBadRequest)
	}
	if u.IsActive {
		return pkg.NewAppError("target is already an employee", nil, http.StatusBadRequest)
	}
	if !u.IsVerified {
		return pkg.NewAppError("target is not verified", nil, http.StatusBadRequest)
	}
	tr.UserID = u.ID
	// Verify if contractor can hire it
	if !f.RoleID.IsGreater(tr.RoleID) {
		return ErrDontHavePermission
	}
	// Verify if the jobs need an establishment or not
	if f.EstablishmentID != 0 {
		tr.EstablishmentID = f.EstablishmentID
	}
	need := es.needEstablishment(tr.RoleID)
	if need && tr.EstablishmentID == 0 {
		return pkg.NewAppError("role needs to have an establishment", nil, http.StatusForbidden)
	}
	if !need && tr.EstablishmentID != 0 {
		return pkg.NewAppError("role need not have an establishment", nil, http.StatusForbidden)
	}
	return nil
}

func (es employeeService) Update(from model.UserRole, target uint, data *model.User) error {
	if target == 0 {
		return ErrUserNotFound
	}
	err := es.canUpdate(from, target)
	if err != nil {
		return err
	}
	data.ID = target
	if err = es.up.Update(data); err != nil {
		return pkg.NewAppError("failed to updated user", err, http.StatusInternalServerError)
	}
	return nil
}

func (es employeeService) Self(userID uint) (model.UserJobs, error) {
	if userID == 0 {
		return model.UserJobs{}, ErrUserNotFound
	}
	jobs, err := es.est.Self(userID)
	if err != nil {
		return model.UserJobs{}, pkg.NewAppError("user not found", err, http.StatusBadRequest)
	}
	return jobs, nil
}

func (es employeeService) Get(f model.UserRole, target uint) (model.UserJobs, error) {
	if f.UserID == 0 || target == 0 {
		return model.UserJobs{}, ErrUserNotFound
	}
	if f.UserID == target {
		return es.est.Self(target)
	}
	if !f.RoleID.IsGreater(model.WAITER) {
		return model.UserJobs{}, ErrDontHavePermission
	}
	jobs, err := es.est.Get(&f, target)
	if err != nil {
		return model.UserJobs{}, pkg.NewAppError("user roles not found", err, http.StatusBadRequest)
	}
	return jobs, nil
}

// WaiterSearch find all waiters in establishment
func (es employeeService) SearchWaiters(eID uint, s *model.Search) ([]model.User, error) {
	us, err := es.est.SearchWaiters(eID, s)
	if err != nil {
		return nil, pkg.NewAppError("no search results", err, http.StatusBadRequest)
	}
	return us, nil
}

// Search find in all employees
func (es employeeService) Search(s *model.SearchEMPL) ([]model.User, error) {
	us, err := es.est.Search(s)
	if err != nil {
		return nil, pkg.NewAppError("no search results", err, http.StatusBadRequest)
	}
	return us, nil
}

// Hire an user by Email and set rol, salary anestablishment, user is hired by a contractor
func (es employeeService) Hire(f model.UserRole, email string, role *model.UserRole) error {
	if err := es.canHire(&f, role, email); err != nil {
		return err
	}
	//TODO CHECK ESTABLISHMENT
	if err := es.est.Hire(role); err != nil {
		return pkg.NewAppError("failed to hire the user", err, http.StatusInternalServerError)
	}
	return nil
}

// HireWaiter Hires a user, assigns him a waiter role and establishes it in the establishment of the contracting party
func (es employeeService) HireWaiter(f model.UserRole, email string, salary float64) error {
	if salary <= 0 {
		return pkg.NewAppError("salary must be a positive number", nil, http.StatusBadRequest)
	}
	tr := model.UserRole{Salary: salary, RoleID: model.WAITER, EstablishmentID: f.EstablishmentID}
	err := es.est.Hire(&tr)
	if err != nil {
		return pkg.NewAppError("failed to hire the user", err, http.StatusInternalServerError)
	}
	return nil
}

func (es employeeService) Fire(from model.UserRole, target uint) error {
	err := es.canUpdate(from, target)
	if err != nil {
		return err
	}
	if err = es.est.Fire(target); err != nil {
		return pkg.NewAppError("failed to fire the user", err, http.StatusInternalServerError)
	}
	return nil
}
