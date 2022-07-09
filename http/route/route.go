package route

import (
	"users-service/http/handler"
	"users-service/http/middleware"

	"github.com/labstack/echo"
)

type route struct {
	mid middleware.Middleware
	uuc handler.UserUC
	euc handler.EMPLuc
	puc handler.ProductUC
	tuc handler.TableUC
}

func NewRouter(mid middleware.Middleware, uuc handler.UserUC, euc handler.EMPLuc, puc handler.ProductUC, tuc handler.TableUC) route {
	return route{
		uuc: uuc,
		mid: mid,
		euc: euc,
		puc: puc,
		tuc: tuc,
	}
}

func (r route) table(e *echo.Echo) {
	g := e.Group("api/v1/table", r.mid.Login)
	g.GET("/:id", r.tuc.Get)
	g.POST("/", r.tuc.CreateIn)
	g.POST("/:id", r.tuc.Create)
	g.DELETE("/", r.tuc.Delete)
	g.DELETE("/:id", r.tuc.DeleteIn)
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

func (r route) employee(e *echo.Echo) {
	g := e.Group("/api/v1/employee", r.mid.Login)
	g.GET("/", r.euc.Get)
	g.GET("/:id", r.euc.GetByID)
	g.POST("/search/", r.euc.Search)
	g.POST("/search/waiter/", r.euc.SearchWaiters)
	g.POST("/hire/:mail", r.euc.Hire)
	g.POST("/hire/waiter/:mail", r.euc.HireWaiter)
	g.PATCH("/fire/:id", r.euc.Fire)
	g.PUT("/:id", r.euc.Update)
}

func (r route) product(e *echo.Echo) {
	g := e.Group("/api/v1/product")
	g.GET("/", r.puc.GetAll)
	// Need login
	g.Use(r.mid.Login)
	g.GET("/:id", r.puc.Get)
	g.GET("/batch/", r.puc.GetInBatch)
	g.POST("/", r.puc.Create)
	g.PUT("/:id", r.puc.Update)
	g.DELETE("/:id", r.puc.Delete)
}

func (r route) Start(e *echo.Echo) {
	e.HTTPErrorHandler = r.mid.Errors
	r.user(e)
	r.employee(e)
	r.product(e)
	r.table(e)
}
