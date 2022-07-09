package storage

import (
	"errors"
	"fmt"
	"users-service/model"

	"gorm.io/gorm"
)

type emplStore struct {
	db *gorm.DB
}

func NewEMPLStore() emplStore {
	return emplStore{db: _db}
}

func (es emplStore) Self(uID uint) (model.UserJobs, error) {
	uj := model.UserJobs{}
	res := es.db.Where("id = ?", uID).
		Select("email", "id", "url", "name", "birth_date", "is_verified").First(&uj.User)
	err := getErrorFromResult(res)
	if err != nil {
		return model.UserJobs{}, fmt.Errorf("find user by id, %w", err)
	}
	res = es.db.Where("user_id = ?", uID).Find(&uj.Jobs)
	err = getErrorFromResult(res)
	if err != nil {
		return model.UserJobs{}, fmt.Errorf("get jobs, %w", err)
	}
	return uj, nil
}

func (es emplStore) Get(from *model.UserRole, target uint) (model.UserJobs, error) {
	uj := model.UserJobs{}
	res := es.db.Where("id = ?", target).
		Select("email", "id", "url", "name", "birth_date", "is_verified").First(&uj.User)
	err := getErrorFromResult(res)
	if err != nil {
		return model.UserJobs{}, fmt.Errorf("find user by id, %w", err)
	}
	res = es.db.Model(&model.UserRole{}).Where("user_id = ?", target).Where("role_id = ?", from.RoleID)
	if from.EstablishmentID != 0 {
		res = res.Where("establishment_id = ?", from.EstablishmentID) //Check business logic
	}
	res = res.Find(&uj.Jobs)
	err = getErrorFromResult(res)
	if err != nil {
		return model.UserJobs{}, fmt.Errorf("get jobs, %w", err)
	}
	return uj, nil
}

func (es emplStore) SearchWaiters(estID uint, s *model.Search) ([]model.User, error) {
	q := s.Query()
	users := []model.User{}
	tx := es.db.Model(&users).Select("users.id", "users.email", "users.name").
		Joins("LEFT JOIN user_roles as r ON users.id = r.user_id").Where("r.establishment_id = ? AND r.role_id = ? AND r.is_active = true", estID, model.WAITER)

	//tx := es.db.Where("establishment_id = ?", estID)
	if q != "" {
		tx = tx.Order(s.Query())
	}
	if s.Limit != 0 {
		tx = tx.Limit(s.Limit)
	}
	if s.Offset != 0 {
		tx = tx.Offset(s.Limit)
	}
	err := getErrorFromResult(tx.Find(&users))
	if err != nil {
		return nil, err
	}
	return users, nil
}

// SELECT u.id, u.name, u.email, r.role_id, r.establishment_id, r.is_active

// FROM users AS u LEFT JOIN user_roles AS r ON u.id = r.user_id
// WHERE r.is_active = false  // OPCIONAL
// AND r.role_id IN (4,8)
// ORDER BY u.id, r.role_id DESC OFFSET 1 LIMIT 8 ;

func (es emplStore) Search(s *model.SearchEMPL) ([]model.User, error) {
	users := []model.User{}
	tx := es.db.Model(&users).Select("users.id", "users.email", "users.name",
		"r.establishment_id", "r.role_id", "r.is_active").
		Joins("LEFT JOIN user_roles as r ON users.id = r.user_id")
	switch s.Status {
	case model.ACTIVE:
		tx = tx.Where("r.is_active = true")
	case model.NOACTVIE:
		tx = tx.Where("r.is_active = false")
	case model.ANY:
	default:
		return nil, errors.New("no status")
	}
	if s.Rols != nil {
		tx.Where("r.role_id IN ?", s.Rols)
	}
	if s.Ests != nil {
		tx.Where("r.establishment_id IN ?", s.Ests)
	}
	q := s.Query()
	if q != "" {
		tx = tx.Order(s.Query())
	}
	if s.Limit != 0 {
		tx = tx.Limit(s.Limit)
	}
	if s.Offset != 0 {
		tx = tx.Offset(s.Limit)
	}
	err := getErrorFromResult(tx.Find(&users))
	if err != nil {
		return nil, err
	}
	return users, nil
}
func (es emplStore) IsActive(uID uint) (bool, error) {
	ur := model.UserRole{}
	res := es.db.Where("user_id =? AND is_active = true", uID).Select("id").First(ur)
	if res.Error != nil {
		return false, res.Error
	}
	if res.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
}

func (es emplStore) FindByEmail(email string) (model.User, error) {
	u := model.User{}
	res := es.db.Select("id", "is_verified").Where("email = ?", email).First(&u)
	err := getErrorFromResult(res)
	if err != nil {
		return model.User{}, err
	}
	r := model.UserRole{}
	res = es.db.Select("role_id", "establishment_id", "is_active").Where("user_id = ? AND is_active = true", u.ID).First(&r)
	err = getErrorFromResult(res)
	if err != nil {
		return model.User{}, err
	}
	u.IsActive = r.IsActive
	u.RoleID = r.RoleID
	u.EstablishmentID = r.EstablishmentID
	// u := model.User{}
	// res := es.db.Select("users.id, r.role_id, r.establishment_id, r.is_active").
	// Joins("LEFT JOIN user_roles AS r ON users.id = r.user_id").Where("r.is_active = true", email).First(&u)
	// err := getErrorFromResult(res)
	// if err != nil {
	// 	return model.User{}, err
	// }
	return u, nil
}

func (es emplStore) Hire(ur *model.UserRole) error {
	ur.Model = model.Model{}
	ur.IsActive = true
	res := es.db.Create(ur)
	return getErrorFromResult(res)
}

func (es emplStore) Fire(uID uint) error {
	res := es.db.Model(&model.UserRole{}).Where("user_id =? AND is_active = true", uID).Update("is_active", false)
	return getErrorFromResult(res)
}
