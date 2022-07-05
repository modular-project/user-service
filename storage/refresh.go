package storage

import (
	"users-service/model"

	"gorm.io/gorm"
)

type refreshStore struct {
	db *gorm.DB
}

func NewRefreshStore() refreshStore {
	return refreshStore{_db}
}

func (rs refreshStore) Create(re *model.Refresh) error {
	return getErrorFromResult(rs.db.Create(re))
}

func (rs refreshStore) Delete(id uint) error {
	return getErrorFromResult(rs.db.Where("id = ?", id).Delete(&model.Refresh{}))
}

func (rs refreshStore) Find(id uint) (model.Refresh, error) {
	m := model.Refresh{}
	res := rs.db.Where("id = ?", id).First(&m)
	return m, getErrorFromResult(res)
}
