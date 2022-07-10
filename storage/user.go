package storage

/*
Implements UserDao interface with gorm
*/

import (
	"errors"
	"fmt"
	"users-service/model"

	"gorm.io/gorm"
)

type userStore struct {
	db *gorm.DB
}

func NewUserStore() userStore {
	return userStore{_db}
}

func (us userStore) Update(user *model.User) error {
	user.UpdatedAt = nil
	res := us.db.Model(user).Select("name", "birth_date", "url", "updated_at").Updates(user)
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (us userStore) Verify(userID uint) error {
	res := us.db.Model(model.User{}).Where("id = ?", userID).Update("is_verified", true)
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (us userStore) Find(id uint) (model.User, error) {
	user := model.User{}
	res := us.db.Model(&user).
		Select("users.id", "users.name", "users.birth_date", "users.url", "users.email", "users.is_verified",
			"r.role_id, r.establishment_id").
		Joins("LEFT JOIN user_roles as r ON users.id = r.user_id AND r.is_active = true").
		Where("users.id = ?", id).
		First(&user)
	if err := getErrorFromResult(res); err != nil {
		return model.User{}, fmt.Errorf("first user: %w", err)
	}
	return user, nil
}

func (us userStore) Create(user *model.User) error {
	user.CreatedAt = nil
	user.UpdatedAt = nil
	user.ID = 0
	user.IsVerified = false
	user.RoleID = 0
	res := us.db.Create(user)
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	return nil
}

func (us userStore) FindByEmail(email *string) (model.User, error) {
	user := model.User{}
	res := us.db.Where("email = ?", email).Select("id", "password").First(&user)
	if err := getErrorFromResult(res); err != nil {
		return model.User{}, fmt.Errorf("first user email = %s: %w", *email, err)
	}
	return user, nil
}

func (us userStore) ChangePassword(id uint, pwd *string) error {
	res := us.db.Model(&model.User{}).Where("id = ?", id).Update("password", pwd)
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("update user: %w", err)
	}
	return nil
}

func (us userStore) IsEmployee(userID uint) (bool, error) {
	m := model.UserRole{}
	res := us.db.Where("user_id = ?", userID).First(&m)
	if res.Error != nil && !errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false, res.Error
	}
	return res.RowsAffected != 0, nil
}
