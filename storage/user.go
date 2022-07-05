package storage

/*
Implements UserDao interface with gorm
*/

import (
	"time"
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
	user.UpdatedAt = time.Time{}
	res := us.db.Model(user).Select("name", "birth_date", "url", "updated_at").Updates(user)
	return getErrorFromResult(res)
}

func (us userStore) Verify(userID uint) error {
	res := us.db.Model(model.User{}).Where("id = ?", userID).Update("is_verified", true)
	return getErrorFromResult(res)
}

func (us userStore) Find(id uint) (model.User, error) {
	user := model.User{}
	res := us.db.Where("id = ?", id).Omit("password").First(&user)
	return user, getErrorFromResult(res)
}

func (us userStore) Create(user *model.User) error {
	user.CreatedAt = time.Time{}
	user.UpdatedAt = time.Time{}
	user.ID = 0
	user.IsVerified = false
	user.RoleID = 0
	res := us.db.Create(user)
	return getErrorFromResult(res)
}

func (us userStore) FindByEmail(email *string) (model.User, error) {
	user := model.User{}
	res := us.db.Where("email = ?", email).Select("id", "password").First(&user)
	return user, getErrorFromResult(res)
}

func (us userStore) ChangePassword(id uint, pwd *string) error {
	res := us.db.Model(&model.User{}).Where("id = ?", id).Update("password", pwd)
	return getErrorFromResult(res)
}

func (us userStore) IsEmployee(userID uint) bool {
	m := model.UserRole{}
	return us.db.Where("user_id = ?", userID).First(&m).RowsAffected != 0
}
