package storage

import (
	"log"
	"users-service/model"

	"gorm.io/gorm"
)

type userSignStore struct {
	db *gorm.DB
}

type kitchenSignStore struct {
	db *gorm.DB
}

func NewUserSignStore() userSignStore {
	return userSignStore{_db}
}

func NewKitchenSignStore() kitchenSignStore {
	return kitchenSignStore{_db}
}

func (us userSignStore) Find(email string) (model.LogIn, error) {
	user := model.User{}
	if us.db == nil {
		log.Fatal("nil db")
	}
	res := us.db.Where("email = ?", email).Select("id", "password").First(&user)
	return model.LogIn{Password: user.Password, ID: user.ID}, getErrorFromResult(res)
}

func (us userSignStore) Create(l *model.LogIn) error {
	res := us.db.Create(&model.User{
		Password: l.Password,
		Email:    l.User,
	})
	return getErrorFromResult(res)
}

func (ks kitchenSignStore) Find(user string) (model.LogIn, error) {
	kit := model.Kitchen{}
	res := ks.db.Where("name = ?", user).Select("id", "password").First(&kit)
	return model.LogIn{Password: kit.Password, ID: kit.ID}, getErrorFromResult(res)
}

func (ks kitchenSignStore) Create(l *model.LogIn) error {
	res := ks.db.Create(&model.Kitchen{
		Password:        l.Password,
		User:            l.User,
		EstablishmentID: l.ID,
	})
	return getErrorFromResult(res)
}
