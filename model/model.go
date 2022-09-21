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

const (
	LOCAL OrderType = iota
	DELIVERY
)

type OrderType int

type RoleID uint

type Model struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	CreatedAt *time.Time     `json:"c_at,omitempty"`
	UpdatedAt *time.Time     `json:"u_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type User struct {
	Model

	Email      string     `gorm:"not null; unique" json:"email,omitempty"`
	Password   string     `gorm:"not null" json:"password,omitempty"`
	URL        *string    `json:"url,omitempty"`
	Name       *string    `json:"name,omitempty"`
	BirthDate  *time.Time `json:"bdate,omitempty"`
	IsVerified bool       `json:"is_verified,omitempty" gorm:"not null"`

	RoleID          RoleID `gorm:"<-:false; -:migration" json:"role_id,omitempty"`
	EstablishmentID uint   `gorm:"<-:false; -:migration" json:"est_id,omitempty"`
	IsActive        bool   `gorm:"<-:false; -:migration" json:"is_active,omitempty"`
	Roles           []Role `gorm:"many2many:user_roles;" json:"-"`
}

type Verification struct {
	UserID    uint
	Code      string         `json:"code,omitempty"`
	ExpiresAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Role struct {
	Model
	Name  string `json:"name,omitempty"`
	Level uint
}

type UserRole struct {
	Model
	UserID          uint    `json:"user_id"`
	RoleID          RoleID  `json:"role_id"`
	EstablishmentID uint    `json:"est_id"`
	IsActive        bool    `json:"is_active"`
	Salary          float64 `json:"salary"`
}

type Kitchen struct {
	Model
	User            string `json:"user,omitempty" gorm:"column:name"`
	Password        string `json:"password,omitempty"`
	EstablishmentID uint
}

type LogIn struct {
	ID       uint
	User     string `json:"user,omitempty"`
	Password string `json:"password,omitempty"`
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
