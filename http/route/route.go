package route

import (
	"users-service/http/handler"
	"users-service/http/middleware"

	"github.com/labstack/echo"
)

type route struct {
	mid middleware.Middleware
	uuc handler.UserUC
}

func NewRouter(mid middleware.Middleware, uuc handler.UserUC) route {
	return route{
		uuc: uuc,
		mid: mid,
	}
}

func (r route) user(e *echo.Echo) {
	g := e.Group("/api/v1/user")
	g.GET("/", r.mid.Login(r.uuc.GetUserData))
	g.POST("/signin/", r.uuc.SignIn)
	g.POST("/signup/", r.uuc.SignUp)
	g.POST("/refresh/", r.uuc.Refresh)
	g.POST("/verify/", r.mid.Login(r.uuc.GenerateVerificationCode))
	g.PUT("/", r.mid.Login(r.uuc.UpdateUserData))
	g.PATCH("/verify/", r.mid.Login(r.uuc.VerifyUser))
	g.PATCH("/password/", r.mid.Login(r.uuc.ChangePassword))
	g.DELETE("/refresh/", r.uuc.SignOut)
}

func (r route) Start(e *echo.Echo) {
	r.user(e)
}
