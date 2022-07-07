package controller

import "users-service/model"

type EMPLStorager interface {
	Get(uint) (model.UserJobs, error)
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
	return es.up.Update(data)
}

func (es employeeService) Get(userID uint) (model.UserJobs, error) {
	if userID == 0 {
		return model.UserJobs{}, ErrUserNotFound
	}
	return es.est.Get(userID)
}

// WaiterSearch find all waiters in establishment
func (es employeeService) SearchWaiters(uID uint, s *model.Search) ([]model.User, error) {
	eID, err := es.can.Equal(uID, model.MANAGER)
	if err != nil {
		return nil, err
	}
	return es.est.SearchWaiters(eID, s)
}

// Search find in all employees
func (es employeeService) Search(uID uint, s *model.SearchEMPL) ([]model.User, error) {
	if err := es.can.Greater(uID, model.MANAGER); err != nil {
		return nil, err
	}
	return es.est.Search(s)
}

// Hire an user by Email and set rol, salary anestablishment, user is hired by a contractor
func (es employeeService) Hire(contractorID uint, email string, role *model.UserRole) error {
	if err := es.can.CanHire(contractorID, email, role); err != nil {
		return err
	}
	return es.est.Hire(role)
}

// HireWaiter Hires a user, assigns him a waiter role and establishes it in the establishment of the contracting party
func (es employeeService) HireWaiter(contractorID uint, email string, salary float64) error {
	e := model.UserRole{Salary: salary, RoleID: model.WAITER}
	if err := es.can.CanHire(contractorID, email, &e); err != nil {
		return err
	}
	return es.est.Hire(&e)
}

func (es employeeService) Fire(from, target uint) error {
	err := es.can.CanUpdate(from, target)
	if err != nil {
		return err
	}
	return es.est.Fire(target)
}
