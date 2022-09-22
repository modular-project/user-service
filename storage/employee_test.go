package storage

import (
	"testing"
	"users-service/model"

	"github.com/stretchr/testify/assert"
)

var _roles = []model.Role{
	{Name: "dueño"},
	{Name: "Admin"},
	{Name: "Gerente"},
	{Name: "Mesero"},
	{Name: "Cocina"},
}

var tables = []struct {
	user model.User
	rols []model.UserRole
	name string
}{
	{
		user: model.User{
			Email: "First@mail.com",
		},
		name: "First",
		rols: []model.UserRole{
			{
				RoleID:   1,
				IsActive: true,
			}, {
				RoleID:          3,
				EstablishmentID: 1,
			},
			{
				RoleID:          3,
				EstablishmentID: 2,
			}, {
				EstablishmentID: 1,
				RoleID:          2,
			},
		},
	},
	{ // Order by name desc, rol desc Where rol In (1,4,3) and status = active
		user: model.User{
			Email: "Second@mail.com",
		},
		rols: []model.UserRole{
			{
				RoleID:          3,
				IsActive:        true,
				EstablishmentID: 1,
			}, {
				RoleID:          2,
				IsActive:        true,
				EstablishmentID: 1,
			}, {
				RoleID:   4,
				IsActive: true,
			}, {
				RoleID:          3,
				EstablishmentID: 2,
			},
		},
	},
}

func TestCleanup(t *testing.T) {
	err := NewDB(TestConfigDB)
	if err != nil {
		t.Fatalf("NewGormDB: %s", err)
	}
	models := []interface{}{model.User{}, model.Role{}, model.UserRole{}}
	err = Drop(models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
}
func TestSearchEmployee(t *testing.T) {
	err := NewDB(TestConfigDB)
	if err != nil {
		t.Fatalf("NewGormDB: %s", err)
	}
	models := []interface{}{model.User{}, model.Role{}}
	t.Cleanup(func() {
		models := []interface{}{model.User{}, model.Role{}, "user_roles"}
		err = _db.Migrator().DropTable(models...)
		if err != nil {
			t.Fatalf("Failed to Create tables: %s", err)
		}
	})

	err = Migrate(nil, models...)
	if err != nil {
		t.Fatalf("Failed to Create tables: %s", err)
	}
	assert := assert.New(t)
	setUsers(assert, t)
	tests := []struct {
		give          model.SearchEMPL
		wantEmployees []model.User
	}{
		{
			// Order by name desc, rol desc Where rol In (1,5,3) and status = active
			give: model.SearchEMPL{
				Search: model.Search{
					OrderBys: []model.OrderBy{
						{
							By:   model.NAME,
							Sort: model.DES,
						},
						{
							By:   model.ROL,
							Sort: model.DES,
						},
					},
					Rols: []uint{1, 4, 3},
				},
			},
			wantEmployees: []model.User{
				{
					Model:           tables[1].user.Model,
					Email:           tables[1].user.Email,
					Name:            tables[1].user.Name,
					RoleID:          4,
					EstablishmentID: 0,
					IsActive:        true,
				}, {
					Model:           tables[1].user.Model,
					Email:           tables[1].user.Email,
					Name:            tables[1].user.Name,
					RoleID:          3,
					EstablishmentID: 1,
					IsActive:        true,
				}, {
					Model:           tables[0].user.Model,
					Email:           tables[0].user.Email,
					Name:            tables[0].user.Name,
					RoleID:          1,
					EstablishmentID: 0,
					IsActive:        true,
				},
			},
		},
	}
	for _, tt := range tests {
		es := NewEMPLStore()
		users, err := es.Search(&tt.give)
		if assert.NoError(err) {
			if assert.Equal(len(tt.wantEmployees), len(users)) {
				for i, u := range users {

					t.Logf("%+v", u)
					assert.Equal(tt.wantEmployees[i].ID, u.ID)
					assert.Equal(tt.wantEmployees[i].Name, u.Name)
					assert.Equal(tt.wantEmployees[i].Email, u.Email)
					assert.Equal(tt.wantEmployees[i].RoleID, u.RoleID, "role")
					assert.Equal(tt.wantEmployees[i].EstablishmentID, u.EstablishmentID, "establishment")
					assert.Equal(tt.wantEmployees[i].IsActive, u.IsActive)
				}
			}
		}
	}
}

var TestConfigDB DBConnection = DBConnection{
	TypeDB:   POSTGRESQL,
	User:     "admin_restaurant",
	Password: "RestAuraNt_pgsql.561965697",
	Host:     "localhost",
	Port:     "5433",
	NameDB:   "testing",
}

// func dropsTables(t *testing.T) {
// 	models := []interface{}{model.User{}, model.Role{}, model.UserRole{}}
// 	err := Drop(models...)
// 	if err != nil {
// 		t.Fatalf("Failed to Create tables: %s", err)
// 	}
// }

func setUsers(assert *assert.Assertions, t *testing.T) {
	_db.CreateInBatches(&_roles, len(_roles))
	for i := range tables {
		if tables[i].name != "" {
			tables[i].user.Name = &tables[i].name
		}
		for y := range tables[i].rols {
			tables[i].rols[y].UserID = uint(i + 1)
		}
		err := _db.Create(&tables[i].user).Error
		if assert.NoError(err) {
			err = _db.CreateInBatches(&tables[i].rols, len(tables[i].rols)).Error
			assert.NoError(err)
		}
	}
}
