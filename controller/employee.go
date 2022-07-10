package controller

import (
	"fmt"
	"net/http"
	"users-service/model"
	"users-service/pkg"
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

type Updater interface {
	Update(*model.User) error
}

type employeeService struct {
	est EMPLStorager
	up  Updater
	can Canner
}

func NewEmployeeService(est EMPLStorager, up Updater, can Canner) employeeService {
	return employeeService{est: est, up: up, can: can}
}

func (es employeeService) Update(from, target uint, data *model.User) error {
	if target == 0 {
		return ErrUserNotFound
	}
	err := es.can.CanUpdate(from, target)
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
		return model.UserJobs{}, pkg.NewAppError("not user not found", err, http.StatusBadRequest)
	}
	return jobs, nil
}

func (es employeeService) Get(from, target uint) (model.UserJobs, error) {
	if from == 0 || target == 0 {
		return model.UserJobs{}, ErrUserNotFound
	}
	if from == target {
		return es.est.Self(target)
	}
	fj, err := es.can.Job(from)
	if err != nil {
		return model.UserJobs{}, err
	}
	if !fj.RoleID.IsGreater(model.WAITER) {
		return model.UserJobs{}, ErrDontHavePermission
	}
	jobs, err := es.est.Get(&fj, target)
	if err != nil {
		return model.UserJobs{}, pkg.NewAppError("user roles not found", err, http.StatusBadRequest)
	}
	return jobs, nil
}

// WaiterSearch find all waiters in establishment
func (es employeeService) SearchWaiters(uID uint, s *model.Search) ([]model.User, error) {
	eID, err := es.can.Equal(uID, model.MANAGER)
	if err != nil {
		return nil, ErrDontHavePermission
	}
	us, err := es.est.SearchWaiters(eID, s)
	if err != nil {
		return nil, pkg.NewAppError("no search results", err, http.StatusBadRequest)
	}
	return us, nil
}

// Search find in all employees
func (es employeeService) Search(uID uint, s *model.SearchEMPL) ([]model.User, error) {
	if err := es.can.Greater(uID, model.MANAGER); err != nil {
		return nil, err
	}
	us, err := es.est.Search(s)
	if err != nil {
		return nil, pkg.NewAppError("no search results", err, http.StatusBadRequest)
	}
	return us, nil
}

// Hire an user by Email and set rol, salary anestablishment, user is hired by a contractor
func (es employeeService) Hire(contractorID uint, email string, role *model.UserRole) error {
	if err := es.can.CanHire(contractorID, email, role); err != nil {
		return err
	}
	//TODO CHECK ESTABLISHMENT
	if err := es.est.Hire(role); err != nil {
		return pkg.NewAppError("could not hire the user", err, http.StatusInternalServerError)
	}
	return nil
}

// HireWaiter Hires a user, assigns him a waiter role and establishes it in the establishment of the contracting party
func (es employeeService) HireWaiter(contractorID uint, email string, salary float64) error {
	e := model.UserRole{Salary: salary, RoleID: model.WAITER}
	err := es.Hire(contractorID, email, &e)
	if err != nil {
		return fmt.Errorf("hire: %w", err)
	}
	return nil
}

func (es employeeService) Fire(from, target uint) error {
	err := es.can.CanUpdate(from, target)
	if err != nil {
		return err
	}
	if err = es.est.Fire(target); err != nil {
		return pkg.NewAppError("could not fire the user", err, http.StatusInternalServerError)
	}
	return nil
}
