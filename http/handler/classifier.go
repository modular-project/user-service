package handler

import (
	"context"
	"mime/multipart"
	"net/http"
	"users-service/pkg"

	"github.com/labstack/echo"
)

type ClassifierServicer interface {
	ClassImg(context.Context, *multipart.FileHeader) (uint32, error)
}

type ClassifierUC struct {
	s ClassifierServicer
}

func NewClassificerServicer(s ClassifierServicer) ClassifierUC {
	return ClassifierUC{s: s}
}

func (cuc ClassifierUC) Classify(c echo.Context) error {
	img, err := c.FormFile("img")
	if err != nil {
		return pkg.NewAppError("Failed to get image", err, http.StatusBadRequest)
	}
	id, err := cuc.s.ClassImg(c.Request().Context(), img)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, responseID(uint64(id)))
}
