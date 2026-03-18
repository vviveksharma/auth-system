package orgmodels

import "time"

// OrgPlan holds plan details returned with org responses.
type OrgPlan struct {
	Name            string  `json:"name"`
	PriceMonthlyUsd float64 `json:"price_monthly_usd"`
}

// CreateOrgOrganizationBody is the organization object inside the CreateOrg response.
type CreateOrgOrganizationBody struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	UserRole    string    `json:"user_role"`
	Plan        OrgPlan   `json:"plan"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateOrgResponseBody is the full response for POST /organizations.
type CreateOrgResponseBody struct {
	Organization CreateOrgOrganizationBody `json:"organization"`
	Message      string                    `json:"message"`
}

// OrgMetadata holds aggregate counts for an organization.
type OrgMetadata struct {
	MemberCount    int     `json:"member_count"`
	ProjectCount   int     `json:"project_count"`
	MonthlyCostUsd float64 `json:"monthly_cost_usd"`
}

// OrgMonthStats holds this-month usage statistics.
type OrgMonthStats struct {
	Requests      int64   `json:"requests"`
	CostUsd       float64 `json:"cost_usd"`
	Tokens        int64   `json:"tokens"`
	Errors        int64   `json:"errors,omitempty"`
	ErrorRate     float64 `json:"error_rate,omitempty"`
	AvgDurationMs int     `json:"avg_duration_ms,omitempty"`
}

// OrgTodayStats holds today's usage statistics.
type OrgTodayStats struct {
	Requests  int64   `json:"requests"`
	CostUsd   float64 `json:"cost_usd"`
	Tokens    int64   `json:"tokens"`
	Errors    int64   `json:"errors"`
	ErrorRate float64 `json:"error_rate"`
}

// OrgPermissions holds what the current user can do in this org.
type OrgPermissions struct {
	CanEdit        bool   `json:"can_edit"`
	CanDelete      bool   `json:"can_delete"`
	CanManageTeam  bool   `json:"can_manage_team"`
	CanViewBilling bool   `json:"can_view_billing"`
	Reason         string `json:"reason,omitempty"`
}

// ListOrgItem is a single organization entry in the list response.
type ListOrgItem struct {
	Id             string        `json:"id"`
	TenantId       string        `json:"tenant_id"`
	Name           string        `json:"name"`
	Slug           string        `json:"slug"`
	Description    string        `json:"description"`
	IconUrl        *string       `json:"icon_url"`
	UserRole       string        `json:"user_role"`
	IsCurrent      bool          `json:"is_current"`
	CurrentBadge   string        `json:"current_badge,omitempty"`
	Metadata       OrgMetadata   `json:"metadata"`
	ThisMonthStats OrgMonthStats `json:"this_month_stats"`
	Plan           OrgPlan       `json:"plan"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

// ListOrgsResponseBody is the full response for GET /organizations.
type ListOrgsResponseBody struct {
	Organizations         []ListOrgItem `json:"organizations"`
	CurrentOrganizationId string        `json:"current_organization_id"`
	TotalCount            int           `json:"total_count"`
}

// OrgMember holds a team member's details.
type OrgMember struct {
	UserId     string    `json:"user_id"`
	Email      string    `json:"email"`
	Name       string    `json:"name"`
	Role       string    `json:"role"`
	JoinedAt   time.Time `json:"joined_at"`
	LastActive time.Time `json:"last_active"`
}

// OrgPlanLimits holds the resource limits for a plan.
type OrgPlanLimits struct {
	Projects         int `json:"projects"`
	TeamMembers      int `json:"team_members"`
	RequestsPerMonth int `json:"requests_per_month"`
}

// OrgPlanUsage holds current usage against plan limits.
type OrgPlanUsage struct {
	Projects          int `json:"projects"`
	TeamMembers       int `json:"team_members"`
	RequestsThisMonth int `json:"requests_this_month"`
}

// OrgPlanDetail holds plan details including limits and usage.
type OrgPlanDetail struct {
	Name            string        `json:"name"`
	PriceMonthlyUsd float64       `json:"price_monthly_usd"`
	Limits          OrgPlanLimits `json:"limits"`
	Usage           OrgPlanUsage  `json:"usage"`
}

// OrgBilling holds billing period details for the org.
type OrgBilling struct {
	CurrentPeriodStart  time.Time `json:"current_period_start"`
	CurrentPeriodEnd    time.Time `json:"current_period_end"`
	NextBillingDate     time.Time `json:"next_billing_date"`
	EstimatedInvoiceUsd float64   `json:"estimated_invoice_usd"`
}

// GetOrgByIdResponseBody is the full response for GET /organizations/{id}.
type GetOrgByIdResponseBody struct {
	Id             string         `json:"id"`
	TenantId       string         `json:"tenant_id"`
	Name           string         `json:"name"`
	Slug           string         `json:"slug"`
	Description    string         `json:"description"`
	IconUrl        *string        `json:"icon_url"`
	OwnerId        string         `json:"owner_id"`
	UserRole       string         `json:"user_role"`
	IsCurrent      bool           `json:"is_current"`
	Metadata       OrgMetadata    `json:"metadata"`
	TodayStats     OrgTodayStats  `json:"today_stats"`
	ThisMonthStats OrgMonthStats  `json:"this_month_stats"`
	Plan           OrgPlanDetail  `json:"plan"`
	Members        []OrgMember    `json:"members"`
	Billing        OrgBilling     `json:"billing"`
	Permissions    OrgPermissions `json:"permissions"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
}

// SwitchOrgResponseBody is returned by POST /organizations/{id}/switch.
type SwitchOrgResponseBody struct {
	OrganizationId         string `json:"organization_id"`
	TenantId               string `json:"tenant_id"`
	Name                   string `json:"name"`
	Switched               bool   `json:"switched"`
	PreviousOrganizationId string `json:"previous_organization_id"`
	Message                string `json:"message"`
}

// UpdateOrgResponseBody is returned by PUT /organizations/{id}.
type UpdateOrgResponseBody struct {
	Id          string    `json:"id"`
	TenantId    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	Message     string    `json:"message"`
}

// DeleteOrgResponseBody is used by the Delete org endpoint.
type DeleteOrgResponseBody struct {
	Message string `json:"message"`
}

// OrgStatsPeriod describes the time window for a stats query.
type OrgStatsPeriod struct {
	Start   string `json:"start"`
	End     string `json:"end"`
	GroupBy string `json:"group_by"`
}

// OrgStatsSummary holds aggregate totals for the stats period.
type OrgStatsSummary struct {
	TotalRequests int64   `json:"total_requests"`
	TotalCostUsd  float64 `json:"total_cost_usd"`
	TotalTokens   int64   `json:"total_tokens"`
	AvgDurationMs int     `json:"avg_duration_ms"`
	ErrorRate     float64 `json:"error_rate"`
}

// OrgStatsDailyItem is one data point in the daily breakdown.
type OrgStatsDailyItem struct {
	Date     string  `json:"date"`
	Requests int64   `json:"requests"`
	CostUsd  float64 `json:"cost_usd"`
	Tokens   int64   `json:"tokens"`
}

// OrgStatsResponseBody is the full response for GET /organizations/{id}/stats.
type OrgStatsResponseBody struct {
	OrganizationId string              `json:"organization_id"`
	TenantId       string              `json:"tenant_id"`
	Period         OrgStatsPeriod      `json:"period"`
	Summary        OrgStatsSummary     `json:"summary"`
	DailyBreakdown []OrgStatsDailyItem `json:"daily_breakdown"`
}
