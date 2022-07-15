package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"users-service/model"
	"users-service/pkg"

	"github.com/gbrlsnchs/jwt"
	"github.com/labstack/echo"
)

type Validater interface {
	Validate(*string) (*jwt.JWT, error)
}

type Permissioner interface {
	UserRole(uint) (model.UserRole, error)
}
type Middleware struct {
	va Validater
	pe Permissioner
}

func NewMiddleware(va Validater, pe Permissioner) Middleware {
	return Middleware{va: va, pe: pe}
}

// Errors handler all errors and checks them to return an response error
func (mid Middleware) Errors(err error, c echo.Context) {
	code, msg := pkg.FindError(err)
	if msg == "" {
		msg = http.StatusText(code)
	}
	res := echo.Map{"message": msg}

	if !c.Response().Committed {
		// Return response with message in JSON
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, res)
		}
		if err != nil {
			// Log new error
			c.Logger().Error(err)
		}
	}

}

func (mid Middleware) authentication(c echo.Context) (jwt.JWT, error) {
	bearer := c.Request().Header.Get("Authorization")
	if bearer == "" {
		return jwt.JWT{}, pkg.NewAppError("user is not logged in", nil, http.StatusUnauthorized)
	}
	splitToken := strings.Split(bearer, "Bearer")
	if len(splitToken) != 2 {
		return jwt.JWT{}, pkg.NewAppError("bearer token not in proper format", nil, http.StatusUnauthorized)
	}
	t := strings.TrimSpace(splitToken[1])
	token, err := mid.va.Validate(&t)
	if err != nil {
		return jwt.JWT{}, fmt.Errorf("mid.va.Validate: %w", err)
	}
	return *token, nil

}

func (mid Middleware) Login(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := mid.authentication(c)
		if err != nil {
			return err
		}
		ty := token.Public()["utp"].(float64)
		if pkg.USER != pkg.UserType(ty) {
			return pkg.NewAppError("it's not a user account", nil, http.StatusUnauthorized)
		}
		c.Set("token", token)
		return next(c)
	}
}

func (mid Middleware) KitchenLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := mid.authentication(c)
		if err != nil {
			return err
		}
		ty := token.Public()["utp"].(float64)
		if pkg.KITCHEN != pkg.UserType(ty) {
			return pkg.NewAppError("it's not a kitchen account", nil, http.StatusUnauthorized)
		}
		c.Set("token", token)
		return next(c)
	}
}
