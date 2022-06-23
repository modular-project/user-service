package model

import (
	"time"

	"gorm.io/gorm"
)

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Model
	Email      string `gorm:"not null"`
	Password   string `gorm:"not null"`
	URL        *string
	Name       *string
	BirthDate  time.Time `json:"bdate"`
	IsVerified bool      `gorm:"not null"`
	RoleID     uint      `gorm:"not null; default 0"`
	Roles      []Role    `gorm:"many2many:user_roles"`
}

type Verification struct {
	UserID    uint
	Code      string
	ExpiresAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Role struct {
	Model
	Name string
}

type UserRole struct {
	Model
	UserID          uint `gorm:"primaryKey"`
	RoleID          uint `gorm:"primaryKey"`
	EstablishmentID uint
	IsActive        bool
	Salary          float64
}

type Kitchen struct {
	Model
	User            string
	Password        string
	EstablishmentID uint
}

type LogIn struct {
	ID       uint
	User     string
	Password string
}

type Refresh struct {
	ID        uint `gorm:"primarykey" json:"id"`
	UserType  int
	Hash      string
	ExpiresAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
