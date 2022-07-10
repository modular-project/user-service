package handler

import (
	"errors"
	"net/http"
	"users-service/pkg"

	"github.com/gbrlsnchs/jwt"
	"github.com/labstack/echo"
)

var (
	ErrUserIDNotFoundInJwt = errors.New("user id not found in jwt")
	ErrUserIDIsNotANumber  = errors.New("user id is not a number")
)

func createRefreshCookie(c echo.Context, refreshToken string) {
	cookie := new(http.Cookie)
	cookie.Name = "refresh"
	cookie.HttpOnly = true
	cookie.Path = "/api/v1/user/refresh/"
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

func getUserIDFromContext(c echo.Context) (uint, error) {
	token := c.Get("token").(jwt.JWT)
	data := token.Public()
	v, ok := data["uid"]
	if !ok {
		return 0, pkg.NewAppError("Fail at get User ID from context", ErrUserIDNotFoundInJwt, http.StatusUnauthorized)
	}
	id, ok := v.(float64)
	if !ok {
		return 0, pkg.NewAppError("Fail at get User ID from context", ErrUserIDIsNotANumber, http.StatusBadRequest)
	}
	return uint(id), nil
}
