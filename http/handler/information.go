package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	pb "github.com/modular-project/protobuffers/information/establishment"
)

func GetEstablishmentByIds(c echo.Context) error {
	req := pb.RequestGetAll{}
	err := c.Bind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, createResponse(fmt.Sprintf("error at bind data: %s", err)))
	}
	return nil
}
