package controller

import (
	"fmt"
	"net/http"
	"testing"
	"users-service/mocks"
	"users-service/model"
	"users-service/pkg"

	"github.com/stretchr/testify/assert"
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
		giveFrom    model.UserRole
		giveEmail   string
		giveRole    model.UserRole
		wantErrCode int
	}{
		{
			giveFrom:  model.UserRole{RoleID: model.OWNER},
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.MANAGER, EstablishmentID: 3, Salary: 140.33},
		}, {
			giveFrom:    model.UserRole{RoleID: model.OWNER},
			giveEmail:   "user@mail.com",
			giveRole:    model.UserRole{RoleID: model.MANAGER, Salary: 140.33},
			wantErrCode: http.StatusForbidden,
		}, {
			giveFrom:    model.UserRole{RoleID: model.OWNER},
			giveEmail:   "user@mail.com",
			giveRole:    model.UserRole{RoleID: model.ADMIN, EstablishmentID: 3, Salary: 140.33},
			wantErrCode: http.StatusForbidden,
		}, {
			giveFrom:    model.UserRole{RoleID: model.OWNER},
			giveEmail:   "manager@mail.com",
			giveRole:    model.UserRole{RoleID: model.WAITER, EstablishmentID: 3, Salary: 140.33},
			wantErrCode: http.StatusBadRequest,
		}, {
			giveFrom:  model.UserRole{RoleID: model.OWNER},
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.ADMIN, Salary: 140.33},
		}, {
			giveFrom:    model.UserRole{RoleID: model.OWNER},
			giveEmail:   "user@mail.com",
			giveRole:    model.UserRole{RoleID: model.OWNER, Salary: 140.33},
			wantErrCode: http.StatusForbidden,
		}, {
			giveFrom:    model.UserRole{},
			giveEmail:   "user@mail.com",
			giveRole:    model.UserRole{RoleID: model.OWNER, Salary: 140.33},
			wantErrCode: http.StatusForbidden,
		}, {
			giveFrom:    model.UserRole{RoleID: model.WAITER, EstablishmentID: 1},
			giveEmail:   "user@mail.com",
			giveRole:    model.UserRole{RoleID: model.WAITER, Salary: 140.33},
			wantErrCode: http.StatusForbidden,
		}, {
			giveFrom:    model.UserRole{RoleID: model.MANAGER, EstablishmentID: 1},
			giveEmail:   "user@mail.com",
			giveRole:    model.UserRole{RoleID: model.MANAGER, Salary: 140.33},
			wantErrCode: http.StatusForbidden,
		}, {
			giveFrom:  model.UserRole{RoleID: model.MANAGER, EstablishmentID: 1},
			giveEmail: "user@mail.com",
			giveRole:  model.UserRole{RoleID: model.WAITER, Salary: 140.33},
		}, {
			giveFrom:    model.UserRole{RoleID: model.OWNER},
			giveEmail:   "noverified@mail.com",
			giveRole:    model.UserRole{RoleID: model.ADMIN, Salary: 140.33},
			wantErrCode: http.StatusBadRequest,
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
		j := mocks.NewPermissionStorage(t)
		p := NewPermission(j)
		uEmail := uByEmail[tt.giveEmail]
		j.On("Find", uEmail.Email).Return(model.User{
			Model:      model.Model{ID: uEmail.ID},
			IsVerified: uEmail.IsVerified,
			IsActive:   uEmail.IsActive,
		}, nil)
		es := NewEmployeeService(nil, nil, p)
		err := es.canHire(&tt.giveFrom, &tt.giveRole, tt.giveEmail)
		code, _ := pkg.FindError(err)
		assert.Equal(t, tt.wantErrCode, code, fmt.Sprintf("Index: %d, uID: %d, email: %s", i, tt.giveFrom.ID, tt.giveEmail))
		j.ExpectedCalls = nil
	}

}
