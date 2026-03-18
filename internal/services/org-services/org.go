package orgservices

import (
	"context"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	orgmodels "github.com/vviveksharma/auth/internal/models/orgModels"
	orgrepo "github.com/vviveksharma/auth/internal/repo/orgRepo"
	"github.com/vviveksharma/auth/models"
)

var slugRegex = regexp.MustCompile(`^[a-z0-9-]+$`)

var planDisplayNames = map[string]string{
	"free":       "Free",
	"pro":        "Pro",
	"enterprise": "Enterprise",
}

var planPricing = map[string]float64{
	"free":       0.0,
	"pro":        99.0,
	"enterprise": 499.0,
}

var planLimits = map[string]orgmodels.OrgPlanLimits{
	"free":       {Projects: 10, TeamMembers: 5, RequestsPerMonth: 100_000},
	"pro":        {Projects: 50, TeamMembers: 20, RequestsPerMonth: 1_000_000},
	"enterprise": {Projects: 10_000, TeamMembers: 100, RequestsPerMonth: 10_000_000},
}

type IOrgServiceInterface interface {
	CreateOrg(ctx context.Context, req *orgmodels.CreateOrgRequestBody) (*orgmodels.CreateOrgResponseBody, error)
	ListOrgs(ctx context.Context) (*orgmodels.ListOrgsResponseBody, error)
	GetOrgById(ctx context.Context, orgId uuid.UUID) (*orgmodels.GetOrgByIdResponseBody, error)
	UpdateOrg(ctx context.Context, orgId uuid.UUID, req *orgmodels.UpdateOrgRequestBody) (*orgmodels.UpdateOrgResponseBody, error)
	SwitchOrg(ctx context.Context, orgId uuid.UUID) (*orgmodels.SwitchOrgResponseBody, error)
	DeleteOrg(ctx context.Context, orgId uuid.UUID) (*orgmodels.DeleteOrgResponseBody, error)
	GetOrgStats(ctx context.Context, orgId uuid.UUID, startDate, endDate, groupBy string) (*orgmodels.OrgStatsResponseBody, error)
}

type OrgService struct {
	OrgRepositoryRepo orgrepo.OrgRepositoryInterface
}

func NewOrgService() (IOrgServiceInterface, error) {
	ser := &OrgService{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (os *OrgService) SetupRepo() error {
	var err error
	organisation, err := orgrepo.NewOrgRepository(db.DB)
	if err != nil {
		return err
	}
	os.OrgRepositoryRepo = organisation
	return err
}

func (os *OrgService) CreateOrg(ctx context.Context, req *orgmodels.CreateOrgRequestBody) (*orgmodels.CreateOrgResponseBody, error) {
	tenantId := ctx.Value("tenant_id").(string)

	if len(req.Name) < 2 || len(req.Name) > 255 {
		return nil, &models.ServiceResponse{Code: 400, Message: "name_required: organization name must be between 2 and 255 characters"}
	}
	if len(req.Slug) < 3 || len(req.Slug) > 100 {
		return nil, &models.ServiceResponse{Code: 400, Message: "slug_invalid: slug must be between 3 and 100 characters"}
	}
	if !slugRegex.MatchString(req.Slug) {
		return nil, &models.ServiceResponse{Code: 400, Message: "slug_invalid: slug must contain only lowercase letters, numbers, and hyphens"}
	}
	if len(req.Description) > 500 {
		return nil, &models.ServiceResponse{Code: 400, Message: "description_invalid: description must not exceed 500 characters"}
	}

	plan := req.Plan
	if plan == "" {
		plan = "free"
	}
	if _, valid := planDisplayNames[plan]; !valid {
		return nil, &models.ServiceResponse{Code: 400, Message: "plan_invalid: plan must be one of free, pro, enterprise"}
	}

	existing, findErr := os.OrgRepositoryRepo.FindByConditons(&models.DBOrganisation{Slug: req.Slug})
	if findErr == nil && existing != nil {
		return nil, &models.ServiceResponse{Code: 409, Message: "slug_taken: an organization with this slug already exists"}
	} else if findErr != nil && findErr.Error() != "record not found" {
		return nil, &models.ServiceResponse{Code: 500, Message: "error while checking slug availability: " + findErr.Error()}
	}

	now := time.Now()
	newOrg := &models.DBOrganisation{
		TenantId:    uuid.MustParse(tenantId),
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		Plan:        plan,
		Status:      "active",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if createErr := os.OrgRepositoryRepo.CreateOrg(newOrg); createErr != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error while creating the organisation: " + createErr.Error()}
	}

	return &orgmodels.CreateOrgResponseBody{
		Organization: orgmodels.CreateOrgOrganizationBody{
			Id:          newOrg.Id.String(),
			Name:        newOrg.Name,
			Slug:        newOrg.Slug,
			Description: newOrg.Description,
			UserRole:    "owner",
			Plan:        orgmodels.OrgPlan{Name: planDisplayNames[plan], PriceMonthlyUsd: planPricing[plan]},
			CreatedAt:   newOrg.CreatedAt,
		},
		Message: "Organization created successfully",
	}, nil
}

func (os *OrgService) ListOrgs(ctx context.Context) (*orgmodels.ListOrgsResponseBody, error) {
	tenantId := ctx.Value("tenant_id").(string)

	orgDetails, totalCount, err := os.OrgRepositoryRepo.ListOrgTenant(uuid.MustParse(tenantId), 1, 1000)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error while fetching organisations: " + err.Error()}
	}

	items := make([]orgmodels.ListOrgItem, 0, len(orgDetails))
	for _, org := range orgDetails {
		plan := org.Plan
		if plan == "" {
			plan = "free"
		}
		var iconUrl *string
		if org.IconUrl != "" {
			iconUrl = &org.IconUrl
		}
		items = append(items, orgmodels.ListOrgItem{
			Id:             org.Id.String(),
			TenantId:       org.TenantId.String(),
			Name:           org.Name,
			Slug:           org.Slug,
			Description:    org.Description,
			IconUrl:        iconUrl,
			UserRole:       "owner",
			IsCurrent:      false,
			Metadata:       orgmodels.OrgMetadata{},
			ThisMonthStats: orgmodels.OrgMonthStats{},
			Plan:           orgmodels.OrgPlan{Name: planDisplayNames[plan], PriceMonthlyUsd: planPricing[plan]},
			CreatedAt:      org.CreatedAt,
			UpdatedAt:      org.UpdatedAt,
		})
	}

	return &orgmodels.ListOrgsResponseBody{
		Organizations:         items,
		CurrentOrganizationId: "",
		TotalCount:            int(totalCount),
	}, nil
}

func (os *OrgService) GetOrgById(ctx context.Context, orgId uuid.UUID) (*orgmodels.GetOrgByIdResponseBody, error) {
	tenantId := ctx.Value("tenant_id").(string)

	orgDetails, err := os.OrgRepositoryRepo.GetOrgById(uuid.MustParse(tenantId), orgId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "organization_not_found: no organisation with this id"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error while fetching org details: " + err.Error()}
	}

	plan := orgDetails.Plan
	if plan == "" {
		plan = "free"
	}
	var iconUrl *string
	if orgDetails.IconUrl != "" {
		iconUrl = &orgDetails.IconUrl
	}

	now := time.Now()
	periodStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	periodEnd := periodStart.AddDate(0, 1, 0)

	return &orgmodels.GetOrgByIdResponseBody{
		Id:             orgDetails.Id.String(),
		TenantId:       orgDetails.TenantId.String(),
		Name:           orgDetails.Name,
		Slug:           orgDetails.Slug,
		Description:    orgDetails.Description,
		IconUrl:        iconUrl,
		OwnerId:        orgDetails.TenantId.String(),
		UserRole:       "owner",
		IsCurrent:      false,
		Metadata:       orgmodels.OrgMetadata{},
		TodayStats:     orgmodels.OrgTodayStats{},
		ThisMonthStats: orgmodels.OrgMonthStats{},
		Plan: orgmodels.OrgPlanDetail{
			Name:            planDisplayNames[plan],
			PriceMonthlyUsd: planPricing[plan],
			Limits:          planLimits[plan],
			Usage:           orgmodels.OrgPlanUsage{},
		},
		Members: []orgmodels.OrgMember{},
		Billing: orgmodels.OrgBilling{
			CurrentPeriodStart:  periodStart,
			CurrentPeriodEnd:    periodEnd,
			NextBillingDate:     periodEnd,
			EstimatedInvoiceUsd: planPricing[plan],
		},
		Permissions: orgmodels.OrgPermissions{
			CanEdit:        true,
			CanDelete:      true,
			CanManageTeam:  true,
			CanViewBilling: true,
		},
		CreatedAt: orgDetails.CreatedAt,
		UpdatedAt: orgDetails.UpdatedAt,
	}, nil
}

func (os *OrgService) UpdateOrg(ctx context.Context, orgId uuid.UUID, req *orgmodels.UpdateOrgRequestBody) (*orgmodels.UpdateOrgResponseBody, error) {
	tenantId := ctx.Value("tenant_id").(string)

	if req.Name != "" && (len(req.Name) < 2 || len(req.Name) > 255) {
		return nil, &models.ServiceResponse{Code: 400, Message: "validation_error: organization name must be between 2 and 255 characters"}
	}
	if len(req.Description) > 500 {
		return nil, &models.ServiceResponse{Code: 400, Message: "validation_error: description must not exceed 500 characters"}
	}

	updatedOrg, err := os.OrgRepositoryRepo.UpdateOrg(uuid.MustParse(tenantId), orgId, &models.DBOrganisation{
		Name:        req.Name,
		Slug:        req.Slug,
		Description: req.Description,
		IconUrl:     req.IconUrl,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "organization_not_found: no organisation with this id"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error while updating the organisation: " + err.Error()}
	}

	return &orgmodels.UpdateOrgResponseBody{
		Id:          updatedOrg.Id.String(),
		TenantId:    updatedOrg.TenantId.String(),
		Name:        updatedOrg.Name,
		Description: updatedOrg.Description,
		UpdatedAt:   updatedOrg.UpdatedAt,
		Message:     "Organization updated successfully",
	}, nil
}

func (os *OrgService) SwitchOrg(ctx context.Context, orgId uuid.UUID) (*orgmodels.SwitchOrgResponseBody, error) {
	tenantId := ctx.Value("tenant_id").(string)

	orgDetails, err := os.OrgRepositoryRepo.GetOrgById(uuid.MustParse(tenantId), orgId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "organization_not_found: no organisation with this id"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error while switching organisation: " + err.Error()}
	}

	return &orgmodels.SwitchOrgResponseBody{
		OrganizationId:         orgDetails.Id.String(),
		TenantId:               orgDetails.TenantId.String(),
		Name:                   orgDetails.Name,
		Switched:               true,
		PreviousOrganizationId: "",
		Message:                "Switched to " + orgDetails.Name,
	}, nil
}

func (os *OrgService) DeleteOrg(ctx context.Context, orgId uuid.UUID) (*orgmodels.DeleteOrgResponseBody, error) {
	tenantId := ctx.Value("tenant_id").(string)

	err := os.OrgRepositoryRepo.DeleteOrg(uuid.MustParse(tenantId), orgId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "organization_not_found: the organisation with this id is not found"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error while deleting the organisation"}
	}

	return &orgmodels.DeleteOrgResponseBody{Message: "Organization deleted successfully"}, nil
}

func (os *OrgService) GetOrgStats(ctx context.Context, orgId uuid.UUID, startDate, endDate, groupBy string) (*orgmodels.OrgStatsResponseBody, error) {
	tenantId := ctx.Value("tenant_id").(string)

	_, err := os.OrgRepositoryRepo.GetOrgById(uuid.MustParse(tenantId), orgId)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "organization_not_found: no organisation with this id"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error while fetching org stats: " + err.Error()}
	}

	if groupBy == "" {
		groupBy = "day"
	}
	if startDate == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	return &orgmodels.OrgStatsResponseBody{
		OrganizationId: orgId.String(),
		TenantId:       tenantId,
		Period:         orgmodels.OrgStatsPeriod{Start: startDate, End: endDate, GroupBy: groupBy},
		Summary:        orgmodels.OrgStatsSummary{},
		DailyBreakdown: []orgmodels.OrgStatsDailyItem{},
	}, nil
}
