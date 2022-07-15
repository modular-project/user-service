package route

import (
	"users-service/http/handler"
	"users-service/http/middleware"
	"users-service/model"

	"github.com/labstack/echo"
)

type route struct {
	mid   middleware.Middleware
	uuc   handler.UserUC
	euc   handler.EMPLuc
	puc   handler.ProductUC
	tuc   handler.TableUC
	estuc handler.ESTDuc
	kuc   handler.KitchenUC
}

func NewRouter(mid middleware.Middleware, uuc handler.UserUC, euc handler.EMPLuc, puc handler.ProductUC, tuc handler.TableUC, estuc handler.ESTDuc,
	kuc handler.KitchenUC) route {
	return route{
		uuc:   uuc,
		mid:   mid,
		euc:   euc,
		puc:   puc,
		tuc:   tuc,
		estuc: estuc,
		kuc:   kuc,
	}
}

func (r route) kitchen(e *echo.Echo) {
	g := e.Group("/api/v1/kitchen")
	g.POST("/signin/", r.kuc.SignIn)
	g.POST("/refresh/", r.kuc.Refresh)
	g.DELETE("/refresh/", r.kuc.SignOut)
	g.Use(r.mid.Equal(model.MANAGER, true))
	g.POST("/signup/", r.kuc.SignUp)
	g.GET("/", r.kuc.GetInESTB)
	g.PUT("/:id", r.kuc.Update)
	g.DELETE("/:id", r.kuc.Delete)
}

func (r route) establishment(e *echo.Echo) {
	g := e.Group("/api/v1/establishment")
	g.GET("/:id", r.estuc.Get)
	g.GET("/", r.estuc.GetInBatch)
	g.Use(r.mid.Greater(model.MANAGER, false))
	g.POST("/", r.estuc.Create)
	g.PUT("/:id", r.estuc.Update)
	g.DELETE("/:id", r.estuc.Delete)
}

func (r route) table(e *echo.Echo) {
	g := e.Group("/api/v1/table")
	g.GET("/:id", r.tuc.Get, r.mid.Login)
	g.POST("/", r.tuc.CreateIn, r.mid.Equal(model.MANAGER, true))
	g.POST("/:id", r.tuc.Create, r.mid.Greater(model.MANAGER, false))
	g.DELETE("/", r.tuc.Delete, r.mid.Equal(model.MANAGER, true))
	g.DELETE("/:id", r.tuc.DeleteIn, r.mid.Greater(model.MANAGER, false))
}

func (r route) user(e *echo.Echo) {
	g := e.Group("/api/v1/user")
	g.POST("/signin/", r.uuc.SignIn)
	g.POST("/signup/", r.uuc.SignUp)
	g.POST("/refresh/", r.uuc.Refresh)
	g.DELETE("/refresh/", r.uuc.SignOut)
	g.Use(r.mid.Login)
	g.GET("/", r.uuc.GetUserData)
	g.POST("/verify/", r.uuc.GenerateVerificationCode)
	g.PUT("/", r.uuc.UpdateUserData)
	g.PATCH("/verify/", r.uuc.VerifyUser)
	g.PATCH("/password/", r.uuc.ChangePassword)
}

func (r route) employee(e *echo.Echo) {
	g := e.Group("/api/v1/employee")
	g.GET("/", r.euc.Get, r.mid.Login)
	g.GET("/:id", r.euc.GetByID, r.mid.Greater(model.WAITER, true))
	g.POST("/search/", r.euc.Search, r.mid.Greater(model.MANAGER, false))
	g.POST("/search/waiter/", r.euc.SearchWaiters, r.mid.Equal(model.MANAGER, true))
	g.POST("/hire/:mail", r.euc.Hire, r.mid.Greater(model.MANAGER, true))
	g.POST("/hire/waiter/:mail", r.euc.HireWaiter, r.mid.Equal(model.MANAGER, true))
	g.PATCH("/fire/:id", r.euc.Fire, r.mid.Greater(model.WAITER, true))
	g.PUT("/:id", r.euc.Update, r.mid.Greater(model.WAITER, true))
}

func (r route) product(e *echo.Echo) {
	g := e.Group("/api/v1/product")
	g.GET("/", r.puc.GetAll)
	g.GET("/:id", r.puc.Get, r.mid.Login)
	g.GET("/batch/", r.puc.GetInBatch, r.mid.Login)
	g.Use(r.mid.Greater(model.MANAGER, false))
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
	r.establishment(e)
	r.kitchen(e)
}
