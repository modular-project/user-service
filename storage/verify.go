package storage

import (
	"users-service/model"

	"gorm.io/gorm"
)

type verifyStore struct {
	db *gorm.DB
}

func NewVerifyStore() verifyStore {
	return verifyStore{db}
}

func (vs verifyStore) Find(userID uint) (model.Verification, error) {
	m := model.Verification{}
	res := vs.db.Where("user_id = ?", userID).First(&m)
	return m, getErrorFromResult(res)
}
func (vs verifyStore) Delete(userID uint) error {
	return getErrorFromResult(vs.db.Where("user_id =?", userID).Delete(&model.Verification{}))
}
func (vs verifyStore) Create(ver *model.Verification) error {
	res := vs.db.Create(ver)
	return getErrorFromResult(res)
}
