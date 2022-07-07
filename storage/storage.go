package storage

import (
	"fmt"
	"sync"
	"users-service/pkg"
)

var once sync.Once

var (
	ErrNotFound = pkg.NotFoundErr("not found")
)

type DBConnection struct {
	TypeDB   DRIVER
	User     string
	Password string
	Port     string
	NameDB   string
	Host     string
}

func NewDB(conn DBConnection) error {
	var err error
	once.Do(func() {
		switch conn.TypeDB {
		case POSTGRESQL:
			err = newPostgresDB(&conn)
		default:
			err = fmt.Errorf("invalid database type")
		}
	})
	return err
}
