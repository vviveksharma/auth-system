package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DBUser struct {
	ID       uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	Role     string    `json:"role"`
}

func (DBUser) TableName() string {
	return "user_tbl"
}

func (*DBUser) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}
