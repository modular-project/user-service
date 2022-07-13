package middleware

import (
	"net/http"
	"users-service/model"
	"users-service/pkg"

	"github.com/labstack/echo"
)

func (mid Middleware) Greater(role model.RoleID, save bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role
			ur, err := mid.userRole(c)
			if err != nil {
				return err
			}
			// Check role
			if !ur.RoleID.IsGreater(role) {
				pkg.NewAppError("you don't have permission", nil, http.StatusForbidden)
			}
			if save {
				c.Set("ur", ur)
			}
			return next(c)
		}
	}
}

func (mid Middleware) Equal(role model.RoleID, save bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Get user role
			ur, err := mid.userRole(c)
			if err != nil {
				return err
			}
			// Check role
			if ur.RoleID != role {
				pkg.NewAppError("you don't have permission", nil, http.StatusForbidden)
			}
			if save {
				c.Set("ur", ur)
			}
			return next(c)
		}
	}
}

// Return an user role if user is login with a user account and have an active rol
func (mid Middleware) userRole(c echo.Context) (model.UserRole, error) {
	token, err := mid.authentication(c)
	if err != nil {
		return model.UserRole{}, err
	}
	pub := token.Public()
	// Check user type
	ty, ok := pub["utp"].(float64)
	if pkg.USER != pkg.UserType(ty) || !ok {
		return model.UserRole{}, pkg.NewAppError("it's not a user account", nil, http.StatusUnauthorized)
	}

	uID, ok := pub["uid"].(float64)
	if !ok {
		return model.UserRole{}, pkg.NewAppError("invalid user id", nil, http.StatusUnauthorized)
	}
	ur, err := mid.pe.UserRole(uint(uID))
	if err != nil {
		return model.UserRole{}, err
	}
	return ur, nil
}
