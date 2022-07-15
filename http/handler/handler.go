package handler

import (
	"errors"
	"net/http"
	"users-service/model"
	"users-service/pkg"

	"github.com/gbrlsnchs/jwt"
	"github.com/labstack/echo"
)

var (
	ErrUserIDNotFoundInJwt = errors.New("user id not found in jwt")
	ErrUserIDIsNotANumber  = errors.New("user id is not a number")
	ErrTokenIsNotAJWT      = errors.New("token is not a jwt")
)

func createRefreshCookie(c echo.Context, refreshToken, path string) {
	cookie := new(http.Cookie)
	cookie.Name = "refresh"
	cookie.HttpOnly = true
	cookie.Path = path
	cookie.MaxAge = 0
	cookie.Value = refreshToken
	c.SetCookie(cookie)
}

func createResponse(message string) map[string]interface{} {
	return map[string]interface{}{
		"message": message,
	}
}

func responseID(id uint64) map[string]interface{} {
	return map[string]interface{}{
		"id": id,
	}
}

func getUserRoleFromContext(c echo.Context) (model.UserRole, error) {
	ur, ok := c.Get("ur").(model.UserRole)
	if !ok {
		return model.UserRole{}, pkg.NewAppError("user don't have a role", nil, http.StatusUnauthorized)
	}
	return ur, nil
}

func getUserIDFromContext(c echo.Context) (uint, error) {
	token, ok := c.Get("token").(jwt.JWT)
	if !ok {
		return 0, pkg.NewAppError("invalid token", ErrTokenIsNotAJWT, http.StatusUnauthorized)
	}
	data := token.Public()
	v, ok := data["uid"]
	if !ok {
		return 0, pkg.NewAppError("fail at get User ID from context", ErrUserIDNotFoundInJwt, http.StatusUnauthorized)
	}
	id, ok := v.(float64)
	if !ok {
		return 0, pkg.NewAppError("fail at get User ID from context", ErrUserIDIsNotANumber, http.StatusBadRequest)
	}
	return uint(id), nil
}
