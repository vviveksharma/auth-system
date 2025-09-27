package models

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type DBUser struct {
	Id        uuid.UUID      `gorm:"primaryKey,column:id"`
	CreatedAt time.Time      `gorm:"column:created_at;not_null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;not_null"`
	TenantId  uuid.UUID      `gorm:"type:uuid;not null"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Password  string         `json:"password"`
	Salt      string         `json:"salt"`
	Status    bool           `json:"status"`
	Roles     pq.StringArray `gorm:"type:text[]" json:"roles"`
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
	Id        uuid.UUID `gorm:"primaryKey,column:id"`
	Role      string    `json:"role"`
	RoleId    uuid.UUID `json:"role_id"`
	TenantId  uuid.UUID `gorm:"type:uuid;not null"`
	RoleType  string    `json:"role_type"`
	Status    bool      `json:"status"`
	CreatedAt time.Time `gorm:"column:created_at;not_null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not_null"`
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
	RoleName  string    `gorm:"type:string;not null"`
	JWTToken  string    `gorm:"type:text;not null"`
	IssuedAt  time.Time `gorm:"autoCreateTime"`
	ExpiresAt time.Time `gorm:"not null"`
	Revoked   bool      `gorm:"default:false;not null"`
	IPAddress string    `gorm:"type:varchar(45)"`
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
	Id        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Salt      string    `json:"salt"`
	Campany   string    `json:"campany"`
	Password  string    `json:"password"`
	Status    string    `json:"status" gorm:"default:'active';index"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
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
	Id             uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantId       uuid.UUID  `gorm:"type:uuid;not null"`
	Name           string     `json:"name"`
	CreatedAt      time.Time  `gorm:"column:created_at;not_null"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;not_null"`
	LastUsedAt     *time.Time `json:"last_used_at"`
	UsageCount     int64      `json:"usage_count" gorm:"default:0"`
	ExpiresAt      time.Time  `gorm:"not null"`
	IsActive       bool       `json:"is_active"`
	ApplicationKey bool       `json:"application_key"`
	RevokedAt      *time.Time `json:"revoked_at"`
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
	Email     string    `json:"email"`
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
	Id          uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	RoleName    string         `json:"role_name"`
	TenantId    uuid.UUID      `gorm:"type:uuid;not null"`
	RoleId      uuid.UUID      `json:"role_id"`
	Permissions string         `gorm:"type:jsonb" json:"permissions"`
	Routes      pq.StringArray `gorm:"type:text[]" json:"routes"`
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
	Id       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserId   uuid.UUID `gorm:"type:uuid;not null;index"`
	TenantId uuid.UUID `gorm:"type:uuid;not null;index"`

	OTPHash string `gorm:"type:varchar(255);not null"`
	OTPType string `gorm:"type:varchar(20);default:'numeric'"`

	ResetToken string    `gorm:"type:varchar(255)"`
	ExpiresAt  time.Time `gorm:"not null;index"`
	IsActive   bool      `gorm:"default:true;index"`

	CreatedAt time.Time  `gorm:"autoCreateTime"`
	UpdatedAt time.Time  `gorm:"autoUpdateTime"`
	UsedAt    *time.Time // Track when OTP was used
}

func (DBResetToken) TableName() string {
	return "db_reset_token"
}

func (*DBResetToken) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}
