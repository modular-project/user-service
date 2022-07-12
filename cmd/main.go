package main

import (
	"fmt"
	"log"
	"os"
	"users-service/adapter/email"
	"users-service/adapter/info"
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
	env := "USER_HOST"
	host, f := os.LookupEnv(env)
	if !f {
		log.Fatalf("environment variable (%s) not found", env)
	}
	env = "USER_PORT"
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
	// env = "ORDER_HOST"
	// oHost, f := os.LookupEnv(env)
	// if !f {
	// 	log.Fatalf("environment variable (%s) not found", env)
	// }
	// env = "ORDER_PORT"
	// oPort, f := os.LookupEnv(env)
	// if !f {
	// 	log.Fatalf("environment variable (%s) not found", env)
	// }
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
	// Migrate tables to DB
	storage.Migrate(
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
	to := authorization.NewToken()
	us := controller.NewUserService(storage.NewUserStore(), storage.NewVerifyStore(), email.NewMail())
	ss := controller.NewSignService(storage.NewRefreshStore(), controller.NewUserValidate(), storage.NewUserSignStore(), to)
	job := storage.NewJobStore()
	per := controller.NewPermission(job)
	es := controller.NewEmployeeService(storage.NewEMPLStore(), storage.NewUserStore(), per)
	// Start GRPC clients
	do := grpc.WithTransportCredentials(insecure.NewCredentials())
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", iHost, iPort), do)
	if err != nil {
		log.Fatalf("fatal at start grpc connection in %s:%s, %v", iHost, iPort, err)
	}
	// Create Adapter dependencies
	ps := info.NewProductService(conn)
	ts := info.NewTableService(conn, per)
	ess := info.NewESTBService(conn)
	// Create Custon Middleware
	mid := mdw.NewMiddleware(to)
	// Create Use Cases
	uUC := handler.NewUserUC(us, ss)
	eUC := handler.NewEMPLUC(es)
	pUC := handler.NewProductUC(ps)
	tUC := handler.NewTableUC(ts)
	estUC := handler.NewESTDuc(ess)
	// Create routes
	r := route.NewRouter(mid, uUC, eUC, pUC, tUC, estUC)
	// Start server
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	r.Start(e)
	err = e.Start(fmt.Sprintf("%s:%s", host, port))
	if err != nil {
		log.Fatalf("%v", err)
	}
}
