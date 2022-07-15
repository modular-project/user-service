package controller

import (
	"fmt"
	"net/http"
	"users-service/model"
	"users-service/pkg"
)

type KitchenStorage interface {
	GetInESTB(uint) ([]model.Kitchen, error)
	Delete(eID, kID uint) error
	Update(eID, kID uint, k *model.Kitchen) error
}

type KitchenService struct {
	kst KitchenStorage
	va  Validater
}

func NewKitchenService(kst KitchenStorage) KitchenService {
	return KitchenService{kst: kst, va: NewKitchenValidate()}
}

func (ks KitchenService) GetInESTB(eID uint) ([]model.Kitchen, error) {
	if eID == 0 {
		return nil, pkg.NewAppError("establishment not found", nil, http.StatusBadRequest)
	}
	kits, err := ks.kst.GetInESTB(eID)
	if err != nil {
		return nil, pkg.NewAppError("failed to get kitchens in establishment", err, http.StatusInternalServerError)
	}
	return kits, nil
}

func (ks KitchenService) Delete(eID, kID uint) error {
	if kID == 0 {
		return pkg.NewAppError("kitchen not found", nil, http.StatusBadRequest)
	}
	err := ks.kst.Delete(eID, kID)
	if err != nil {
		return pkg.NewAppError("failed to delete kitchen", err, http.StatusInternalServerError)
	}
	return nil
}

func (ks KitchenService) Update(eID, kID uint, k *model.Kitchen) error {
	if kID == 0 {
		return pkg.NewAppError("kitchen not found", nil, http.StatusBadRequest)
	}
	l := model.LogIn{User: k.User, Password: k.Password}
	if err := ks.va.Validate(&l); err != nil {
		return fmt.Errorf("validate: %w", err)
	}
	err := ks.kst.Update(eID, kID, k)
	if err != nil {
		return pkg.NewAppError("failed to delete kitchen", err, http.StatusInternalServerError)
	}
	return nil
}
