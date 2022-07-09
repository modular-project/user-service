package storage

import (
	"fmt"
	"users-service/model"

	"gorm.io/gorm"
)

type jobStore struct {
	db *gorm.DB
}

func NewJobStore() jobStore {
	return jobStore{db: _db}
}

func (js jobStore) Find(email string) (model.User, error) {
	user := model.User{}
	res := js.db.Model(&user).Where("email = ?", email).
		Select("users.id", "users.is_verified", "r.role_id", "r.establishment_id", "r.is_active").
		Joins("LEFT JOIN user_roles as r ON r.user_id = users.id").Last(&user)
	if err := getErrorFromResult(res); err != nil {
		return model.User{}, fmt.Errorf("user %w", err)
	}
	return user, nil
}

func (js jobStore) Job(uID uint) (model.UserRole, error) {
	r := model.UserRole{}
	res := js.db.Where("user_id = ? and is_active = true", uID).Select("user_id, role_id, establishment_id, is_active").First(&r)
	return r, getErrorFromResult(res)
}

func (js jobStore) IsEmployee(userID uint) (bool, error) {
	m := model.UserRole{}
	res := js.db.Select("id").Where("user_id = ?", userID).First(&m)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

func (js jobStore) IsActive(userID uint) (bool, error) {
	m := model.UserRole{}
	res := js.db.Select("id").Where("user_id = ? AND is_active = true", userID).First(&m)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

func (js jobStore) IsVerified(userID uint) (bool, error) {
	m := model.User{}
	res := js.db.Select("is_verified").Where("user_id = ?", userID).First(&m)
	return m.IsVerified, getErrorFromResult(res)
}
