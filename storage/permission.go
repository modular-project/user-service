package storage

import (
	"fmt"
	"users-service/model"

	"gorm.io/gorm"
)

type permissionStore struct {
	db *gorm.DB
}

func NewPermissionStore() permissionStore {
	return permissionStore{db: _db}
}

func (js permissionStore) Kitchen(kID uint) (uint, error) {
	k := model.Kitchen{Model: model.Model{ID: kID}}
	res := js.db.First(&k)
	if err := getErrorFromResult(res); err != nil {
		return 0, fmt.Errorf("first kitchen: %w", err)
	}
	return k.EstablishmentID, nil
}

func (js permissionStore) Find(email string) (model.User, error) {
	user := model.User{}
	res := js.db.Model(&user).Where("email = ?", email).
		Select("users.id", "users.is_verified", "r.role_id", "r.establishment_id", "r.is_active").
		Joins("LEFT JOIN user_roles as r ON r.user_id = users.id").Order("r.is_active DESC").First(&user)
	if err := getErrorFromResult(res); err != nil {
		return model.User{}, fmt.Errorf("last user %w", err)
	}
	return user, nil
}

func (js permissionStore) UserRole(uID uint) (model.UserRole, error) {
	r := model.UserRole{}
	res := js.db.Where("user_id = ? and is_active = true", uID).Select("user_id, role_id, establishment_id, is_active").First(&r)
	if err := getErrorFromResult(res); err != nil {
		return model.UserRole{}, fmt.Errorf("first user: %w", err)
	}
	return r, nil
}

func (js permissionStore) IsEmployee(userID uint) (bool, error) {
	m := model.UserRole{}
	res := js.db.Select("id").Where("user_id = ?", userID).First(&m)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

func (js permissionStore) IsActive(userID uint) (bool, error) {
	m := model.UserRole{}
	res := js.db.Select("id").Where("user_id = ? AND is_active = true", userID).First(&m)
	if res.Error != nil {
		return false, res.Error
	}
	return res.RowsAffected == 1, nil
}

func (js permissionStore) IsVerified(userID uint) (bool, error) {
	m := model.User{}
	res := js.db.Select("is_verified").Where("user_id = ?", userID).First(&m)
	if err := getErrorFromResult(res); err != nil {
		return false, fmt.Errorf("first user: %w", err)
	}
	return m.IsVerified, nil
}
