package storage

import (
	"fmt"
	"users-service/model"

	"gorm.io/gorm"
)

type kitchenStore struct {
	db *gorm.DB
}
type KitchenStorage interface {
	GetInESTB(uint) ([]model.Kitchen, error)
	Delete(eID, kID uint) error
	Update(eID, kID uint, k *model.Kitchen) error
}

func NewKitchenStore() kitchenStore {
	return kitchenStore{db: _db}
}

func (ks kitchenStore) GetInESTB(eID uint) ([]model.Kitchen, error) {
	var kits []model.Kitchen
	res := ks.db.Where("establishment_id = ?", eID).Find(&kits)
	if err := getErrorFromResult(res); err != nil {
		return nil, fmt.Errorf("find kitchens by estbID: %w", err)
	}
	return kits, nil
}
func (ks kitchenStore) Delete(eID, kID uint) error {
	res := ks.db.Where("id = ?", kID).Where("establishment_id = ?", eID).Delete(&model.Kitchen{})
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("delete kitchen: %w", err)
	}
	return nil
}
func (ks kitchenStore) Update(eID, kID uint, k *model.Kitchen) error {
	k.UpdatedAt = nil
	res := ks.db.Select("name, password, updated_at").Where("id = ?", kID).Where("establishment_id = ?", eID).Updates(k)
	if err := getErrorFromResult(res); err != nil {
		return fmt.Errorf("update kitchen: %w", err)
	}
	return nil
}
