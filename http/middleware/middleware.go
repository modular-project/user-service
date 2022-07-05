package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"users-service/pkg"

	"github.com/gbrlsnchs/jwt"
	"github.com/labstack/echo"
)

var (
	ErrUserNotLoggedIn      = errors.New("user is not logged in")
	ErrBearerTokenNotFormat = errors.New("bearer token not in proper format")
)

type Validater interface {
	Validate(*string) (*jwt.JWT, error)
}

type Middleware struct {
	va Validater
}

func NewMiddleware(va Validater) Middleware {
	return Middleware{va: va}
}

func (mid Middleware) authentication(c echo.Context) (jwt.JWT, error) {
	if bearer := c.Request().Header.Get("Authorization"); bearer != "" {
		splitToken := strings.Split(bearer, "Bearer")
		if len(splitToken) != 2 {
			return jwt.JWT{}, ErrBearerTokenNotFormat
		}
		t := strings.TrimSpace(splitToken[1])
		token, err := mid.va.Validate(&t)
		if err != nil {
			return jwt.JWT{}, err
		}
		return *token, nil
	}
	return jwt.JWT{}, ErrUserNotLoggedIn
}

func createResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"msg": message,
	}
}

func (mid Middleware) Login(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := mid.authentication(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, createResponse(fmt.Sprintf("user not logged in: %s", err)))
		}
		ty := token.Public()["utp"].(float64)
		if pkg.USER != pkg.UserType(ty) {
			return c.JSON(http.StatusUnauthorized, createResponse("not a user"))
		}
		c.Set("token", token)
		return next(c)
	}
}

func (mid Middleware) KitchenLogin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, err := mid.authentication(c)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, createResponse(fmt.Sprintf("user not logged in: %s", err)))
		}
		ty := token.Public()["utp"].(float64)
		if pkg.KITCHEN != pkg.UserType(ty) {
			return c.JSON(http.StatusUnauthorized, createResponse("not a kitchen"))
		}
		c.Set("token", token)
		return next(c)
	}
}
