package projectmodels

import (
	"time"

	"github.com/vviveksharma/auth/internal/pagination"
)

// ── Shared stats sub-types ────────────────────────────────────────────────────

type ProjectRequestCounts struct {
	Total      int64   `json:"total"`
	Successful int64   `json:"successful"`
	Failed     int64   `json:"failed"`
	ErrorRate  float64 `json:"error_rate"`
}

type ProjectTokenCounts struct {
	Total         int64 `json:"total"`
	Input         int64 `json:"input,omitempty"`
	Output        int64 `json:"output,omitempty"`
	AvgPerRequest int64 `json:"avg_per_request,omitempty"`
}

type ProjectPerformance struct {
	AvgDurationMs int `json:"avg_duration_ms"`
	P50DurationMs int `json:"p50_duration_ms,omitempty"`
	P95DurationMs int `json:"p95_duration_ms,omitempty"`
	P99DurationMs int `json:"p99_duration_ms,omitempty"`
}

type ProjectTodayStats struct {
	Date        string               `json:"date"`
	Requests    ProjectRequestCounts `json:"requests"`
	CostUsd     float64              `json:"cost_usd"`
	Tokens      ProjectTokenCounts   `json:"tokens"`
	Performance ProjectPerformance   `json:"performance"`
}

type ProjectMonthStats struct {
	Year        int                  `json:"year"`
	Month       int                  `json:"month"`
	Requests    ProjectRequestCounts `json:"requests"`
	CostUsd     float64              `json:"cost_usd"`
	Tokens      ProjectTokenCounts   `json:"tokens"`
	Performance ProjectPerformance   `json:"performance"`
}

type ProjectListStats struct {
	Requests  int64   `json:"requests"`
	CostUsd   float64 `json:"cost_usd"`
	Tokens    int64   `json:"tokens"`
	Errors    int64   `json:"errors"`
	ErrorRate float64 `json:"error_rate"`
}

type ProjectApiKeysSummary struct {
	ActiveCount int `json:"active_count"`
	TotalCount  int `json:"total_count"`
}

// ── 1. List projects ──────────────────────────────────────────────────────────

type ListProjectItem struct {
	Id          string                `json:"id"`
	Name        string                `json:"name"`
	Description string                `json:"description"`
	Environment string                `json:"environment"`
	TodayStats  ProjectListStats      `json:"today_stats"`
	MonthStats  ProjectListStats      `json:"month_stats"`
	ApiKeys     ProjectApiKeysSummary `json:"api_keys"`
	CreatedAt   time.Time             `json:"created_at"`
	UpdatedAt   *time.Time            `json:"updated_at"`
}

type ListPagination = pagination.PaginationMeta

type ListProjectsResponseBody struct {
	OrganizationId string                    `json:"organization_id"`
	Projects       []ListProjectItem         `json:"projects"`
	Pagination     pagination.PaginationMeta `json:"pagination"`
}

// ── 2. Project details ────────────────────────────────────────────────────────

type ProjectDetail struct {
	Id             string     `json:"id"`
	OrganizationId string     `json:"organization_id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	Environment    string     `json:"environment"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type ProjectProviderSimple struct {
	Provider      string   `json:"provider"`
	ProviderLabel string   `json:"provider_label"`
	Models        []string `json:"models"`
	Requests      int64    `json:"requests"`
	CostUsd       float64  `json:"cost_usd"`
	CostShare     float64  `json:"cost_share"`
	Tokens        int64    `json:"tokens"`
}

type AlertConfig struct {
	Enabled             bool    `json:"enabled"`
	ThresholdUsd        float64 `json:"threshold_usd,omitempty"`
	ThresholdPercentage float64 `json:"threshold_percentage,omitempty"`
	Period              string  `json:"period"`
}

type AlertConfiguration struct {
	CostAlert      AlertConfig `json:"cost_alert"`
	ErrorRateAlert AlertConfig `json:"error_rate_alert"`
}

type ApiKeyEntry struct {
	Id        string    `json:"id"`
	KeyPrefix string    `json:"key_prefix"`
	Name      string    `json:"name"`
	Status    string    `json:"status"`
	LastUsed  time.Time `json:"last_used"`
	CreatedAt time.Time `json:"created_at"`
}

type ApiKeysSummaryDetail struct {
	ActiveCount int           `json:"active_count"`
	TotalCount  int           `json:"total_count"`
	Keys        []ApiKeyEntry `json:"keys"`
}

type GetProjectDetailsResponseBody struct {
	Project            ProjectDetail           `json:"project"`
	TodayStats         ProjectTodayStats       `json:"today_stats"`
	MonthStats         ProjectMonthStats       `json:"month_stats"`
	ProvidersUsage     []ProjectProviderSimple `json:"providers_usage"`
	AlertConfiguration AlertConfiguration      `json:"alert_configuration"`
	ApiKeysSummary     ApiKeysSummaryDetail    `json:"api_keys_summary"`
}

// ── 3. Providers breakdown ────────────────────────────────────────────────────

type ProviderModelBreakdown struct {
	Model      string  `json:"model"`
	ModelLabel string  `json:"model_label"`
	Requests   int64   `json:"requests"`
	CostUsd    float64 `json:"cost_usd"`
	Tokens     int64   `json:"tokens"`
}

type ProviderBreakdownItem struct {
	Provider      string                   `json:"provider"`
	ProviderLabel string                   `json:"provider_label"`
	Models        []ProviderModelBreakdown `json:"models"`
	TotalRequests int64                    `json:"total_requests"`
	TotalCostUsd  float64                  `json:"total_cost_usd"`
	CostShare     float64                  `json:"cost_share"`
	TotalTokens   int64                    `json:"total_tokens"`
}

type ProvidersBreakdownPeriod struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type ProvidersBreakdownResponseBody struct {
	ProjectId     string                   `json:"project_id"`
	Period        ProvidersBreakdownPeriod `json:"period"`
	Providers     []ProviderBreakdownItem  `json:"providers"`
	TotalRequests int64                    `json:"total_requests"`
	TotalCostUsd  float64                  `json:"total_cost_usd"`
	TotalTokens   int64                    `json:"total_tokens"`
}

// ── 4. Create project ─────────────────────────────────────────────────────────

type CreatedProject struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Environment string     `json:"environment"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type CreatedApiKey struct {
	Id        string    `json:"id"`
	Key       string    `json:"key"`
	KeyPrefix string    `json:"key_prefix"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Warning   string    `json:"warning"`
}

type CreateProjectResponseBody struct {
	Project CreatedProject `json:"project"`
	ApiKey  *CreatedApiKey `json:"api_key,omitempty"`
}

// ── 5. Update project ─────────────────────────────────────────────────────────

type UpdateProjectResponseBody struct {
	Id          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Environment string     `json:"environment"`
	UpdatedAt   *time.Time `json:"updated_at"`
	Message     string     `json:"message"`
}

// ── 6. Delete project ─────────────────────────────────────────────────────────

type DeleteSideEffects struct {
	ApiKeysRevoked    int  `json:"api_keys_revoked"`
	AlertsDeleted     int  `json:"alerts_deleted"`
	UsageDataArchived bool `json:"usage_data_archived"`
}

type DeleteProjectResponseBody struct {
	Id          string            `json:"id"`
	Name        string            `json:"name"`
	Deleted     bool              `json:"deleted"`
	DeletedAt   time.Time         `json:"deleted_at"`
	Message     string            `json:"message"`
	SideEffects DeleteSideEffects `json:"side_effects"`
}

// ── 7. Project errors ─────────────────────────────────────────────────────────

type ProjectErrorEntry struct {
	Id           string    `json:"id"`
	Timestamp    time.Time `json:"timestamp"`
	Provider     string    `json:"provider"`
	Model        string    `json:"model"`
	Endpoint     string    `json:"endpoint"`
	StatusCode   int       `json:"status_code"`
	ErrorType    string    `json:"error_type"`
	ErrorMessage string    `json:"error_message"`
	DurationMs   int       `json:"duration_ms"`
	CostUsd      float64   `json:"cost_usd"`
}

type ProjectErrorsResponseBody struct {
	ProjectId   string              `json:"project_id"`
	Date        string              `json:"date"`
	TotalErrors int                 `json:"total_errors"`
	ErrorRate   float64             `json:"error_rate"`
	Errors      []ProjectErrorEntry `json:"errors"`
}
