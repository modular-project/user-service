package model

import (
	"time"

	"gorm.io/gorm"
)

const (
	USER RoleID = iota
	OWNER
	ADMIN
	MANAGER
	WAITER
)

type RoleID uint

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt time.Time      `json:"c_at,omitempty"`
	UpdatedAt time.Time      `json:"u_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Model

	Email      string    `gorm:"not null; unique"`
	Password   string    `gorm:"not null" json:",omitempty"`
	URL        *string   `json:",omitempty"`
	Name       *string   `json:",omitempty"`
	BirthDate  time.Time `json:"bdate,omitempty"`
	IsVerified bool      `gorm:"not null"`

	RoleID          RoleID `gorm:"<-:false; -:migration" json:",omitempty"`
	EstablishmentID uint   `gorm:"<-:false; -:migration" json:",omitempty"`
	IsActive        bool   `gorm:"<-:false; -:migration" json:",omitempty"`
	Roles           []Role `gorm:"many2many:user_roles;"`
}

type Verification struct {
	UserID    uint
	Code      string
	ExpiresAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Role struct {
	Model
	Name  string
	Level uint
}

type UserRole struct {
	Model
	UserID          uint
	RoleID          RoleID
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

type UserJobs struct {
	User User       `json:"user"`
	Jobs []UserRole `json:"jobs"`
}

func (r RoleID) IsGreater(target RoleID) bool {
	if r == target {
		return false
	}
	return r < target && r != USER || target == USER
}
