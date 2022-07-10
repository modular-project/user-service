package storage

import (
	"fmt"
	"users-service/model"

	"gorm.io/gorm"
)

type verifyStore struct {
	db *gorm.DB
}

func NewVerifyStore() verifyStore {
	return verifyStore{_db}
}

func (vs verifyStore) Find(userID uint) (model.Verification, error) {
	m := model.Verification{}
	res := vs.db.Where("user_id = ?", userID).First(&m)
	if err := getErrorFromResult(res); err != nil {
		return model.Verification{}, fmt.Errorf("first verification: %w", err)
	}
	return m, nil
}
func (vs verifyStore) Delete(userID uint) error {
	res := vs.db.Where("user_id =?", userID).Delete(&model.Verification{})
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("delete verification: %w", err)
	}
	return nil
}
func (vs verifyStore) Create(ver *model.Verification) error {
	res := vs.db.Create(ver)
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("create verification: %w", err)
	}
	return nil
}
