package storage

import (
	"errors"
	"fmt"
	"log"
	"users-service/pkg"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DRIVER string

const (
	POSTGRESQL DRIVER = "POSTGRES"
	TESTING    DRIVER = "TESTING"
)

var (
	db *gorm.DB
)

func getErrorFromResult(tx *gorm.DB) error {
	if errors.Is(tx.Error, gorm.ErrRecordNotFound) {
		return pkg.ErrNoRowsAffected
	}
	if tx.Error != nil {
		return tx.Error
	}
	if tx.RowsAffected == 0 {
		return pkg.ErrNoRowsAffected
	}
	return nil
}

func newPostgresDB(u *DBConnection) error {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		u.Host, u.User, u.Password, u.NameDB, u.Port)
	db, err = gorm.Open(postgres.Open(dsn))
	if db == nil {
		log.Fatalf("nil db at open db")
	}
	if err != nil {
		return fmt.Errorf("open postgres: %w", err)
	}
	log.Println("connected to postgres")
	return nil
}

func Drop(tables ...interface{}) error {
	return db.Migrator().DropTable(tables...)
}

func Migrate(tables ...interface{}) error {
	return db.AutoMigrate(tables...)
}
