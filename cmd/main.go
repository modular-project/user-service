package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"users-service/adapter/classifier"
	"users-service/adapter/email"
	"users-service/adapter/gdrive"
	"users-service/adapter/info"
	"users-service/adapter/order"
	"users-service/authorization"
	"users-service/controller"
	"users-service/http/handler"
	"users-service/http/route"
	"users-service/model"
	"users-service/storage"

	mdw "users-service/http/middleware"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func newDBConnection() storage.DBConnection {
	env := "USER_DB_USER"
	u, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "USER_DB_PWD"
	pass, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "USER_DB_NAME"
	name, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "USER_DB_HOST"
	host, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "USER_DB_PORT"
	port, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	return storage.DBConnection{
		TypeDB:   storage.POSTGRESQL,
		User:     u,
		Password: pass,
		Port:     port,
		NameDB:   name,
		Host:     host,
	}
}

func main() {
	// Load Credentials

	env := "USER_PORT"
	port, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_HOST"
	iHost, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "INFO_PORT"
	iPort, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "ORDER_HOST"
	oHost, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "ORDER_PORT"
	oPort, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "ADDR_HOST"
	aHost, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "ADDR_PORT"
	aPort, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "IA_HOST"
	cHost, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "IA_PORT"
	cPort, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	err := authorization.LoadCertificates(authorization.RSA512)
	if err != nil {
		log.Fatalf("no se pudo cargar los certificados: %v", err)
	}
	err = email.LoadMail()
	if err != nil {
		log.Fatalf("no se pudo conectar al servicio de mensajeria: %v", err)
	}
	err = storage.NewDB(newDBConnection())
	if err != nil {
		log.Fatalf("NewGormDB(%+v): %s", storage.DBConnection{}, err)
	}
	to := authorization.NewToken()
	rs := storage.NewRefreshStore()
	ss := controller.NewSignService(rs, controller.NewUserValidate(), storage.NewUserSignStore(), to)
	// Migrate tables to DB
	err = storage.Migrate(
		ss,
		&model.User{},
		&model.Role{},
		&model.UserRole{},
		&model.Refresh{},
		&model.Kitchen{},
		&model.Verification{},
	)
	if err != nil {
		log.Fatalf("no se logro realizar la migracion de las tablas: %v", err)
	}
	// Create Dependencies

	us := controller.NewUserService(storage.NewUserStore(), storage.NewVerifyStore(), email.NewMail())
	kss := controller.NewSignService(rs, controller.NewKitchenValidate(), storage.NewKitchenSignStore(), to)
	job := storage.NewPermissionStore()
	per := controller.NewPermission(job)
	es := controller.NewEmployeeService(storage.NewEMPLStore(), storage.NewUserStore(), per)
	ks := controller.NewKitchenService(storage.NewKitchenStore())
	// Start GRPC clients
	do := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", iHost, iPort), do)
	if err != nil {
		log.Fatalf("fatal at start grpc connection to Information Server in %s:%s, %v", iHost, iPort, err)
	}
	oConn, err := grpc.Dial(fmt.Sprintf("%s:%s", oHost, oPort), do)
	if err != nil {
		log.Fatalf("fatal at start grpc connection to Order Server in %s:%s, %v", oHost, oPort, err)
	}
	aConn, err := grpc.Dial(fmt.Sprintf("%s:%s", aHost, aPort), do)
	if err != nil {
		log.Fatalf("fatal at start grpc connection to Address Server in %s:%s, %v", aHost, aPort, err)
	}
	cConn, err := grpc.Dial(fmt.Sprintf("%s:%s", cHost, cPort), do)
	if err != nil {
		log.Fatalf("fatal at start grpc connection to Address Server in %s:%s, %v", cHost, cPort, err)
	}
	// Create Adapter dependencies
	ps := info.NewProductService(conn)
	ts := info.NewTableService(conn)
	ess := info.NewESTBService(conn)
	ios := info.NewInfoOrderService(conn)
	ads := info.NewAddressService(aConn)
	os := order.NewOrderService(oConn, ios)
	oss := order.NewOrderStatusService(oConn, ads, ess)
	imgs := gdrive.NewService()
	cis := classifier.NewClassifierService(cConn)
	// Create Custon Middleware
	mid := mdw.NewMiddleware(to, per)
	// Create Use Cases
	uUC := handler.NewUserUC(us, ss)
	eUC := handler.NewEMPLUC(es)
	pUC := handler.NewProductUC(ps)
	tUC := handler.NewTableUC(ts)
	estUC := handler.NewESTDuc(ess, ads)
	kUC := handler.NewKitchenUC(kss, ks)
	oUC := handler.NewOrderUC(os)
	osUC := handler.NewOrderStatusUC(oss)
	deUC := handler.NewAddDeliveryUC(ads)
	iUC := handler.NewImageUC(imgs)
	cuc := handler.NewClassificerServicer(cis)
	// Create routes
	//TODO ADD NEAREST TO ORDER SERVICE
	r := route.NewRouter(mid, uUC, eUC, pUC, tUC, estUC, kUC, oUC, osUC, deUC, iUC, cuc)
	// Start server
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
	}))
	// Health Check without logs
	e.GET("/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello World") })
	e.GET("/api/v1/", func(c echo.Context) error { return c.String(http.StatusOK, "Hello Api v1") })
	e.Use(middleware.Logger())
	r.Start(e)
	err = e.Start(fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
