package handler

import (
	"mime/multipart"
	"net/http"
	"users-service/adapter/gdrive"
	"users-service/pkg"

	"github.com/labstack/echo"
)

type ImageServicer interface {
	SaveImg(*multipart.FileHeader, string, string) (string, error)
}

type ImageUC struct {
	is ImageServicer
}

func NewImageUC(is ImageServicer) ImageUC {
	return ImageUC{is: is}
}

func (iuc ImageUC) UploadUser(c echo.Context) error {
	img, err := c.FormFile("img")
	if err != nil {
		return pkg.NewAppError("Failed to get image", err, http.StatusBadRequest)
	}
	name := c.FormValue("name")
	link, err := iuc.is.SaveImg(img, name, gdrive.UFOLDER)
	if err != nil {
		return pkg.NewAppError("Could not save image", err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, createResponse(link))
}

func (iuc ImageUC) UploadProduct(c echo.Context) error {
	img, err := c.FormFile("img")
	if err != nil {
		return pkg.NewAppError("Failed to get image", err, http.StatusBadRequest)
	}
	name := c.FormValue("name")
	link, err := iuc.is.SaveImg(img, name, gdrive.PFOLDER)
	if err != nil {
		return pkg.NewAppError("Could not save image", err, http.StatusBadRequest)
	}
	return c.JSON(http.StatusOK, createResponse(link))
}
