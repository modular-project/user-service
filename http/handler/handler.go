package handler

import (
	"errors"
	"net/http"

	"github.com/gbrlsnchs/jwt"
	"github.com/labstack/echo"
)

var (
	ErrRefreshTokenNotFound = errors.New("refresh token not found")
	ErrUserIdNotFoundInJwt  = errors.New("user id not found in jwt")
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
		"msg": message,
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
		return 0, ErrUserIdNotFoundInJwt
	}
	return uint(v.(float64)), nil
}
