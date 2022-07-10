package storage

import (
	"errors"
	"fmt"
	"log"
	"users-service/model"
	"users-service/pkg"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DRIVER string

const (
	POSTGRESQL DRIVER = "POSTGRES"
	TESTING    DRIVER = "TESTING"
)

var (
	_db *gorm.DB
)

func getErrorFromResult(tx *gorm.DB) error {
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return pkg.ErrNoRowsAffected
	}
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return pkg.ErrNoRowsAffected
	}
	return nil
}

func newPostgresDB(u *DBConnection) error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		u.Host, u.User, u.Password, u.NameDB, u.Port)
	_db, err = gorm.Open(postgres.Open(dsn))
	if _db == nil {
		log.Fatalf("nil db at open db")
	}
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}
	log.Println("connected to postgres")
	return nil
}

func setRoles() error {
	roles := []model.Role{
		{
			Model: model.Model{ID: uint(model.OWNER)},
			Name:  "Owner",
		}, {
			Model: model.Model{ID: uint(model.ADMIN)},
			Name:  "Admin",
		}, {
			Model: model.Model{ID: uint(model.MANAGER)},
			Name:  "Manager",
		}, {
			Model: model.Model{ID: uint(model.WAITER)},
			Name:  "Waiter",
		},
	}
	isRoles := []model.Role{}
	res := _db.Find(&isRoles)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 0 {
		return nil
	}
	res = _db.CreateInBatches(&roles, len(roles))
	log.Println("roles was created")
	return getErrorFromResult(res)
}

func setOwner() error {
	u := model.User{}
	res := _db.Where("email = ?", "owner@mail.com").First(&u)
	err := getErrorFromResult(res)
	if err != nil {
		return err
	}
	if u.IsVerified {
		return nil
	}
	res = _db.Model(&model.User{}).Where("email = ?", "owner@mail.com").Update("is_verified", true)
	err = getErrorFromResult(res)
	if err != nil {
		return err
	}
	res = _db.Create(&model.UserRole{
		UserID:   u.ID,
		RoleID:   model.OWNER,
		Salary:   100,
		IsActive: true,
	})
	log.Println("owner is set")
	return getErrorFromResult(res)
}

func Drop(tables ...interface{}) error {
	return _db.Migrator().DropTable(tables...)
}

func Migrate(tables ...interface{}) error {
	err := _db.SetupJoinTable(&model.User{}, "Roles", &model.UserRole{})
	if err != nil {
		return fmt.Errorf("fail at setup join table :%w", err)
	}
	err = _db.AutoMigrate(tables...)
	if err != nil {
		return err
	}
	err = setRoles()
	if err != nil {
		return err
	}
	return setOwner()
}

// func DB() *gorm.DB {
// 	return _db
// }
