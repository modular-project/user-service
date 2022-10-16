package storage

import (
	"errors"
	"fmt"
	"log"
	"os"
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

type creater interface {
	SignUp(l *model.LogIn) error
}

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
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable ",
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
		}, {
			Model: model.Model{ID: uint(model.CHEF)},
			Name:  "Chef",
		},
	}
	isRoles := []model.Role{}
	res := _db.Find(&isRoles)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected != 0 {
		if res.RowsAffected < 5 {
			_db.Create(&model.Role{Model: model.Model{ID: uint(model.CHEF)}, Name: "Chef"})
		}
		return nil
	}

	res = _db.CreateInBatches(&roles, len(roles))
	log.Println("roles was created")
	return getErrorFromResult(res)
}

func setOwner(c creater) error {

	u := model.User{}
	res := _db.Where("id = 1").First(&u)
	log.Println(res.Error != nil && res.Error != gorm.ErrRecordNotFound)
	if res.Error != nil && res.Error != gorm.ErrRecordNotFound {
		return res.Error
	}
	if res.RowsAffected == 1 {
		return nil
	}
	email, ok := os.LookupEnv("USER_OWNER_MAIL")
	if !ok {
		return fmt.Errorf("owner email not found")
	}
	pwd, ok := os.LookupEnv("USER_OWNER_PWD")
	if !ok {
		return fmt.Errorf("owner email not found")
	}
	err := c.SignUp(&model.LogIn{User: email, Password: pwd})
	if err != nil {
		return fmt.Errorf("failed to create owner: %w", err)
	}
	res = _db.Create(&model.UserRole{
		UserID:   1,
		RoleID:   model.OWNER,
		Salary:   0,
		IsActive: true,
	})
	log.Println("owner is set")
	return getErrorFromResult(res)
}

func Drop(tables ...interface{}) error {
	return _db.Migrator().DropTable(tables...)
}

func Migrate(c creater, tables ...interface{}) error {
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
	return setOwner(c)
}

// func DB() *gorm.DB {
// 	return _db
// }
