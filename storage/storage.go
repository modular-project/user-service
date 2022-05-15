package storage

import (
	"fmt"
	"log"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DRIVER string

const (
	POSTGRESQL DRIVER = "POSTGRES"
	TESTING    DRIVER = "TESTING"
)

var (
	db   *gorm.DB
	once sync.Once
)

type dbUser struct {
	TypeDB   DRIVER
	User     string
	Password string
	Port     string
	NameDB   string
	Host     string
}

func New(driver DRIVER) {
	once.Do(func() {
		u := loadData()
		switch u.TypeDB {
		case POSTGRESQL:
			newPostgresDB(&u)
		case TESTING:
			newTestingDB(&u)
		}
	})
}

func newTestingDB(u *dbUser) {
	var err error
	fmt.Print(u)
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", u.User, u.Password, u.Host, u.Port, "testing")
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}

	fmt.Println("conectado a Testing")
}

func newPostgresDB(u *dbUser) {
	var err error
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", u.User, u.Password, u.Host, u.Port, u.NameDB)
	db, err = gorm.Open(postgres.Open(dsn))
	if err != nil {
		log.Fatalf("no se pudo abrir la base de datos: %v", err)
	}

	fmt.Println("conectado a postgres")
}

// DB return a unique instance of db
func DB() *gorm.DB {
	return db
}

func getEnv(env string) (string, error) {
	s, f := os.LookupEnv(env)
	if !f {
		return "", fmt.Errorf("environment variable (%s) not found", env)
	}
	return s, nil
}

func loadData() dbUser {
	typeDb, err := getEnv("RGE_TYPE")
	if err != nil {
		log.Fatalf(err.Error())
	}
	user, err := getEnv("RGE_USER")
	if err != nil {
		log.Fatalf(err.Error())
	}
	password, err := getEnv("RGE_PASSWORD")
	if err != nil {
		log.Fatalf(err.Error())
	}
	port, err := getEnv("RGE_PORT")
	if err != nil {
		log.Fatalf(err.Error())
	}
	name, err := getEnv("RGE_NAME_DB")
	if err != nil {
		log.Fatalf(err.Error())
	}
	host, err := getEnv("RGE_HOST")
	if err != nil {
		log.Fatalf(err.Error())
	}
	return dbUser{DRIVER(typeDb), user, password, port, name, host}
}
