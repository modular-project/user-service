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
	ErrUserNotLoggedIn      = pkg.UnauthorizedErr("user is not logged in")
	ErrBearerTokenNotFormat = pkg.UnauthorizedErr("bearer token not in proper format")
)

type Validater interface {
	Validate(*string) (*jwt.JWT, error)
}

type Middleware struct {
	va Validater
}

type IsForbiddener interface {
	IsForbidden()
}

type IsBadder interface {
	IsBad()
}

type IsUnauthorizeder interface {
	IsUnauthorized()
}

type IsNotFounder interface {
	IsNotFound()
}

func NewMiddleware(va Validater) Middleware {
	return Middleware{va: va}
}

func findError(err error) (int, error) {
	temp := err
	for temp != nil {
		switch temp.(type) {
		case IsBadder:
			return http.StatusBadRequest, temp
		case IsUnauthorizeder:
			return http.StatusUnauthorized, temp
		case IsForbiddener:
			return http.StatusForbidden, temp

		case IsNotFounder:
			return http.StatusNotFound, temp
		}
		temp = errors.Unwrap(temp)
	}
	return http.StatusInternalServerError, nil
}

// Errors handler all errors and checks them to return an response error
func (mid Middleware) Errors(err error, c echo.Context) {
	var msg interface{}

	code, gErr := findError(err)
	if gErr != nil {
		msg = gErr.Error()
	} else {
		msg = http.StatusText(code)
	}

	if _, ok := msg.(string); ok {
		msg = map[string]interface{}{"message": msg}
	}
	// Log error
	//c.Logger().Error(err)

	if !c.Response().Committed {
		// Return response with message in JSON
		if c.Request().Method == http.MethodHead {
			err = c.NoContent(code)
		} else {
			err = c.JSON(code, msg)
		}
		if err != nil {
			// Log new error
			c.Logger().Error(err)
		}
	}

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
