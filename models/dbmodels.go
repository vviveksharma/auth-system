package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DBUser struct {
	Id       uuid.UUID      `gorm:"primaryKey,column:id"`
	TenantId uuid.UUID      `gorm:"type:uuid;not null"`
	Name     string         `json:"name"`
	Email    string         `json:"email"`
	Password string         `json:"password"`
	Salt     string         `json:"salt"`
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
	TenantId uuid.UUID `gorm:"type:uuid;not null"`
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
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantId  uuid.UUID `gorm:"type:uuid;not null"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	RoleId    uuid.UUID `gorm:"type:uuid;not null"`
	JWTToken  string    `gorm:"type:text;not null"`
	IssuedAt  time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false;not null"`
}

func (DBLogin) TableName() string {
	return "login_tbl"
}

func (*DBLogin) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBTenant struct {
	Id       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name     string    `json:"name"`
	Email    string    `json:"email"`
	Campany  string    `json:"campany"`
	Password string    `json:"password"`
}

func (DBTenant) TableName() string {
	return "tenant_tbl"
}

func (*DBTenant) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBToken struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantId  uuid.UUID `gorm:"type:uuid;not null"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `gorm:"not null"`
	IsActive  bool      `json:"is_active"`
}

func (DBToken) TableName() string {
	return "token_tbl"
}

func (*DBToken) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBTenantLogin struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantId  uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	IsActive  bool      `json:"is_active"`
	IPAddress string    `json:"ip_address"`
}

func (DBTenantLogin) TableName() string {
	return "tenant_login_tbl"
}

func (*DBTenantLogin) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBRouteRole struct {
	Id       uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantId uuid.UUID      `gorm:"type:uuid;not null"`
	RoleId   uuid.UUID      `json:"role_id"`
	Route    pq.StringArray `gorm:"type:text[]" json:"routes"`
}

func (DBRouteRole) TableName() string {
	return "route_role_tbl"
}

func (*DBRouteRole) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBResetToken struct {
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId    uuid.UUID `gorm:"type:uuid;not null"`
	TenantId  uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt time.Time `gorm:"not null"`
	IsActive  bool      `json:"is_active"`
}

func (DBResetToken) TableName() string {
	return "db_reset_token"
}

func (*DBResetToken) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}
