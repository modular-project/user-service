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

func getUserIDFromContext(c echo.Context) (uint, error) {
	token := c.Get("token").(jwt.JWT)
	data := token.Public()
	v, ok := data["uid"]
	if !ok {
		return 0, ErrUserIdNotFoundInJwt
	}
	return uint(v.(float64)), nil
}

// func createTokenCookie(r *http.Request, w *http.ResponseWriter, token string) error {
// 	s, err := cookiestorage.DB().Get(r, "sessions")
// 	if err != nil {
// 		return err
// 	}
// 	s.Options = &sessions.Options{
// 		Path:     "/",
// 		Domain:   "",
// 		MaxAge:   time.Now().Minute() * 30,
// 		HttpOnly: true,
// 	}
// 	s.Values["token"] = token
// 	err = s.Save(r, *w)
// 	return err
// }
