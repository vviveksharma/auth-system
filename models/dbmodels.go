package models

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DBUser struct {
	Id       uuid.UUID      `gorm:"primaryKey,column:id"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Roles    pq.StringArray `gorm:"type:text[]" json:"roles"`
}

func (DBUser) TableName() string {
	return "user_tbl"
}

func (*DBUser) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBRoles struct {
	Id       uuid.UUID `gorm:"primaryKey,column:id"`
	Role     string    `json:"role"`
	RoleId   uuid.UUID `json:"role_id"`
	RoleType string    `json:"role_type"`
}

func (DBRoles) TableName() string {
	return "role_tbl"
}

func (*DBRoles) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBLogin struct {
	Id     uuid.UUID `gorm:"primaryKey,column:id"`
	UserId uuid.UUID `json:"user_id"`
	RoleId uuid.UUID `json:"role_id"`
	JWT    string    `json:"jwt"`
}

func (DBLogin) TableName() string {
	return "login_tbl"
}

func (*DBLogin) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}
