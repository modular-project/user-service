package controller_test

import (
	"errors"
	"testing"
	"users-service/controller"
	"users-service/mocks"
	"users-service/model"
)

func TestCanHire(t *testing.T) {
	users := []model.User{
		{
			Model:      model.Model{ID: 1},
			Email:      "owner@mail.com",
			RoleID:     model.OWNER,
			IsActive:   true,
			IsVerified: true,
		}, {
			Model:           model.Model{ID: 2},
			Email:           "manager@mail.com",
			RoleID:          model.MANAGER,
			EstablishmentID: 1,
			IsVerified:      true,
			IsActive:        true,
		}, {
			Model:      model.Model{ID: 3},
			Email:      "user@mail.com",
			IsVerified: true,
		}, {
			Model:      model.Model{ID: 4},
			Email:      "noverified@mail.com",
			IsVerified: false,
		}, {
			Model:           model.Model{ID: 5},
			Email:           "waiter@mail.com",
			RoleID:          model.WAITER,
			EstablishmentID: 1,
			IsActive:        true,
			IsVerified:      true,
		},
	}

	tests := []struct {
		giveID    uint
		giveEmail string
		giveRole  model.UserRole
		wantErr   error
	}{
		{
			giveID:    1,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.MANAGER, EstablishmentID: 3, Salary: 140.33},
		}, {
			giveID:    1,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.MANAGER, Salary: 140.33},
			wantErr:   controller.ErrEstablishNecesary,
		}, {
			giveID:    1,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.ADMIN, EstablishmentID: 3, Salary: 140.33},
			wantErr:   controller.ErrCannotBeAssigned,
		}, {
			giveID:    1,
			giveEmail: "manager@mail.com",
			giveRole:  model.UserRole{RoleID: model.WAITER, EstablishmentID: 3, Salary: 140.33},
			wantErr:   controller.ErrAlreadyEmployee,
		}, {
			giveID:    1,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.ADMIN, Salary: 140.33},
		}, {
			giveID:    1,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.OWNER, Salary: 140.33},
			wantErr:   controller.ErrUnauthorizedUser,
		}, {
			giveID:    3,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.OWNER, Salary: 140.33},
			wantErr:   controller.ErrUnauthorizedUser,
		}, {
			giveID:    5,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.WAITER, Salary: 140.33},
			wantErr:   controller.ErrUnauthorizedUser,
		}, {
			giveID:    2,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.MANAGER, Salary: 140.33},
			wantErr:   controller.ErrUnauthorizedUser,
		}, {
			giveID:    2,
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.WAITER, Salary: 140.33},
			wantErr:   nil,
		}, {
			giveID:    1,
			giveEmail: "noverified@mail.com",
			giveRole:  model.UserRole{RoleID: model.ADMIN, Salary: 140.33},
			wantErr:   controller.ErrUserIsNotVerified,
		},
	}
	uByID := make(map[uint]*model.User, len(users))
	uByEmail := make(map[string]*model.User)
	for i := range users {
		u := users[i]
		uByID[u.ID] = &u
		uByEmail[u.Email] = &u
	}
	for i, tt := range tests {
		j := mocks.NewJobStorager(t)
		p := controller.NewPermission(j)
		uID := uByID[tt.giveID]
		uEmail := uByEmail[tt.giveEmail]
		j.On("Job", uID.ID).Return(model.UserRole{
			UserID:          uID.ID,
			RoleID:          uID.RoleID,
			EstablishmentID: uID.EstablishmentID,
			IsActive:        uID.IsActive,
		}, nil)
		j.On("Find", uEmail.Email).Return(model.User{
			Model:      model.Model{ID: uEmail.ID},
			IsVerified: uEmail.IsVerified,
			IsActive:   uEmail.IsActive,
		}, nil)
		err := p.CanHire(tt.giveID, tt.giveEmail, &tt.giveRole)
		if !errors.Is(err, tt.wantErr) {
			t.Logf("%d - %d: %s", i, tt.giveID, tt.giveEmail)
			t.Errorf("got error: %s, want error: %s", err, tt.wantErr)
		} else {
			j.ExpectedCalls = nil
		}
	}

}
