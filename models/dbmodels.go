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
	OrgId     uuid.UUID      `gorm:"type:uuid;not null"`
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBRoles struct {
	Id          uuid.UUID `gorm:"primaryKey,column:id"`
	Role        string    `json:"role"`
	DisplayName string    `json:"display_name"`
	Description string    `json:"description"`
	RoleId      uuid.UUID `json:"role_id"`
	TenantId    uuid.UUID `gorm:"type:uuid;not null"`
	RoleType    string    `json:"role_type"`
	Status      bool      `json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;not_null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not_null"`
}

func (DBRoles) TableName() string {
	return "role_tbl"
}

func (*DBRoles) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
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
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBMessage struct {
	Id            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserEmail     string    `gorm:"type:varchar(255);not null" json:"user_email"`
	TenantId      uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	CurrentRole   string    `gorm:"type:varchar(100);not null" json:"current_role"`
	RequestedRole string    `gorm:"type:varchar(100);not null" json:"requested_role"`
	Status        string    `gorm:"type:varchar(50);default:'pending'" json:"status"`
	RequestAt     time.Time `gorm:"autoCreateTime" json:"request_at"`
	Action        bool      `gorm:"default:false" json:"action"`
}

func (DBMessage) TableName() string {
	return "message_tbl"
}

func (*DBMessage) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBOrganisation struct {
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	TenantId    uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	Name        string    `json:"name"`
	Slug        string    `gorm:"uniqueIndex" json:"slug"`
	Description string    `gorm:"type:text" json:"description"`
	IconUrl     string    `gorm:"type:text" json:"icon_url"`
	Plan        string    `gorm:"type:varchar(50);default:'free'" json:"plan"`
	Status      string    `gorm:"type:varchar(50);default:'active'" json:"status"`
	CreatedAt   time.Time `gorm:"column:created_at;not_null"`
	UpdatedAt   time.Time `gorm:"column:updated_at;not_null"`
}

func (DBOrganisation) TableName() string {
	return "organisation_tbl"
}

func (*DBOrganisation) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBResetCreds struct {
	Id        uuid.UUID  `gorm:"type:uuid;primaryKey"`
	TenantId  uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	UserId    uuid.UUID  `gorm:"type:uuid;not null" json:"user_id"`
	Active    bool       `json:"active"`
	CodeHash  string     `gorm:"type:text;not null;index"`
	Salt      string     `gorm:"type:text;not null;index"`
	CreatedAt time.Time  `json:"created_at"`
	UsedAt    *time.Time `json:"used_at"`
}

func (DBResetCreds) TableName() string {
	return "reset_creds_tbl"
}

func (*DBResetCreds) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New()
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBProject struct {
	Id          uuid.UUID  `gorm:"type:uuid;primaryKey"`
	TenantId    uuid.UUID  `gorm:"type:uuid;not null" json:"tenant_id"`
	OrgId       uuid.UUID  `gorm:"type:uuid;not null" json:"org_id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Environment string     `json:"environment"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

func (DBProject) TableName() string {
	return "project_tbl"
}

func (*DBProject) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBProjectDailyStats struct {
	Id                 uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProjectId          uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	OrganizationId     uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	TenantId           uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	Date               time.Time `gorm:"type:date;not null" json:"date"`
	TotalRequests      int64     `gorm:"default:0" json:"total_requests"`
	SuccessfulRequests int64     `gorm:"default:0" json:"successful_requests"`
	FailedRequests     int64     `gorm:"default:0" json:"failed_requests"`
	TotalTokens        int64     `gorm:"default:0" json:"total_tokens"`
	TotalCostUsd       float64   `gorm:"type:decimal(12,6);default:0" json:"total_cost_usd"`
	AvgDurationMs      int       `gorm:"default:0" json:"avg_duration_ms"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DBProjectDailyStats) TableName() string {
	return "project_daily_stats"
}

func (*DBProjectDailyStats) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBProjectMonthlyStats struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	ProjectId      uuid.UUID `gorm:"type:uuid;not null;index" json:"project_id"`
	OrganizationId uuid.UUID `gorm:"type:uuid;not null;index" json:"organization_id"`
	TenantId       uuid.UUID `gorm:"type:uuid;not null" json:"tenant_id"`
	Year           int       `gorm:"not null" json:"year"`
	Month          int       `gorm:"not null" json:"month"`
	TotalRequests  int64     `gorm:"default:0" json:"total_requests"`
	TotalTokens    int64     `gorm:"default:0" json:"total_tokens"`
	TotalCostUsd   float64   `gorm:"type:decimal(12,6);default:0" json:"total_cost_usd"`
	ApiKeysCount   int       `gorm:"default:0" json:"api_keys_count"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DBProjectMonthlyStats) TableName() string {
	return "project_monthly_stats"
}

func (*DBProjectMonthlyStats) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}

type DBProviderDailyStats struct {
	Id             uuid.UUID  `gorm:"type:uuid;primaryKey"`
	OrganizationId uuid.UUID  `gorm:"type:uuid;not null;index" json:"organization_id"`
	ProjectId      *uuid.UUID `gorm:"type:uuid;index" json:"project_id"`
	Provider       string     `gorm:"type:varchar(50);not null;index" json:"provider"`
	Date           time.Time  `gorm:"type:date;not null" json:"date"`
	TotalRequests  int64      `gorm:"default:0" json:"total_requests"`
	TotalTokens    int64      `gorm:"default:0" json:"total_tokens"`
	TotalCostUsd   float64    `gorm:"type:decimal(12,6);default:0" json:"total_cost_usd"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
}

func (DBProviderDailyStats) TableName() string {
	return "provider_daily_stats"
}

func (*DBProviderDailyStats) BeforeCreate(tx *gorm.DB) error {
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}
	tx.Statement.SetColumn("Id", uuid)
	return nil
}
