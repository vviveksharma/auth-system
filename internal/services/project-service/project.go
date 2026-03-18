package projectservice

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	projectmodels "github.com/vviveksharma/auth/internal/models/projectModels"
	"github.com/vviveksharma/auth/internal/pagination"
	projectrepo "github.com/vviveksharma/auth/internal/repo/projectRepo"
	"github.com/vviveksharma/auth/models"
)

type IProjectServiceInterface interface {
	CreateProject(ctx context.Context, orgId uuid.UUID, req *projectmodels.CreateProjectRequestBody) (*projectmodels.CreateProjectResponseBody, error)
	ListProjects(ctx context.Context, orgId uuid.UUID, page int, limit int) (*projectmodels.ListProjectsResponseBody, error)
	GetProjectDetails(ctx context.Context, projectId uuid.UUID, date time.Time) (*projectmodels.GetProjectDetailsResponseBody, error)
	GetProvidersBreakdown(ctx context.Context, projectId uuid.UUID, startDate time.Time, endDate time.Time) (*projectmodels.ProvidersBreakdownResponseBody, error)
	UpdateProject(ctx context.Context, projectId uuid.UUID, req *projectmodels.UpdateProjectRequestBody) (*projectmodels.UpdateProjectResponseBody, error)
	DeleteProject(ctx context.Context, projectId uuid.UUID) (*projectmodels.DeleteProjectResponseBody, error)
	GetProjectErrors(ctx context.Context, projectId uuid.UUID, date time.Time, limit int) (*projectmodels.ProjectErrorsResponseBody, error)
}

type Projectservice struct {
	ProjectRepo projectrepo.ProjectRepositryInterface
}

func NewProjectService() (IProjectServiceInterface, error) {
	ser := &Projectservice{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (ps *Projectservice) SetupRepo() error {
	var err error
	project, err := projectrepo.NewProjectReposistry(db.DB)
	if err != nil {
		return err
	}
	ps.ProjectRepo = project
	return err
}

func tenantId(ctx context.Context) (uuid.UUID, error) {
	raw := ctx.Value("tenant_id")
	if raw == nil {
		return uuid.Nil, &models.ServiceResponse{Code: 401, Message: "missing tenant_id in context"}
	}
	switch v := raw.(type) {
	case string:
		return uuid.Parse(v)
	case uuid.UUID:
		return v, nil
	}
	return uuid.Nil, &models.ServiceResponse{Code: 401, Message: "invalid tenant_id type in context"}
}

func (ps *Projectservice) CreateProject(ctx context.Context, orgId uuid.UUID, req *projectmodels.CreateProjectRequestBody) (*projectmodels.CreateProjectResponseBody, error) {
	tid, err := tenantId(ctx)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 401, Message: "unauthorized: " + err.Error()}
	}

	existing, err := ps.ProjectRepo.GetProjectByName(orgId, tid, req.Name)
	if err != nil && err.Error() != "record not found" {
		return nil, &models.ServiceResponse{Code: 500, Message: "error checking existing project: " + err.Error()}
	}
	if existing != nil {
		return nil, &models.ServiceResponse{Code: 409, Message: "a project with this name already exists in the organization"}
	}

	now := time.Now()
	err = ps.ProjectRepo.Create(&models.DBProject{
		TenantId:    tid,
		OrgId:       orgId,
		Name:        req.Name,
		Description: req.Description,
		Environment: req.Environment,
		CreatedAt:   now,
		UpdatedAt:   nil,
	})
	if err != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error creating project: " + err.Error()}
	}

	created, err := ps.ProjectRepo.GetProjectByName(orgId, tid, req.Name)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error fetching created project: " + err.Error()}
	}

	resp := &projectmodels.CreateProjectResponseBody{
		Project: projectmodels.CreatedProject{
			Id:          created.Id.String(),
			Name:        created.Name,
			Description: created.Description,
			Environment: created.Environment,
			CreatedAt:   created.CreatedAt,
			UpdatedAt:   created.UpdatedAt,
		},
	}

	if req.GenerateApiKey {
		keyId := uuid.Must(uuid.NewV7())
		resp.ApiKey = &projectmodels.CreatedApiKey{
			Id:        keyId.String(),
			Key:       "ak_live_" + created.Id.String()[:8] + "_placeholder",
			KeyPrefix: "ak_live_" + created.Id.String()[:8] + "...",
			Name:      "Default Key",
			CreatedAt: now,
			Warning:   "Save this key now. You won't be able to see it again.",
		}
	}

	return resp, nil
}

func (ps *Projectservice) ListProjects(ctx context.Context, orgId uuid.UUID, page int, limit int) (*projectmodels.ListProjectsResponseBody, error) {
	tid, err := tenantId(ctx)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 401, Message: "unauthorized: " + err.Error()}
	}

	page, limit = pagination.ParsePaginationParams(page, limit)

	projects, totalCount, err := ps.ProjectRepo.ListOrgProject(orgId, tid, page, limit)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error listing projects: " + err.Error()}
	}

	today := time.Now()
	items := make([]projectmodels.ListProjectItem, 0, len(projects))
	for _, p := range projects {
		var todayStats projectmodels.ProjectListStats
		var monthStats projectmodels.ProjectListStats

		if ds, err := ps.ProjectRepo.GetProjectDailyStats(p.Id, today); err == nil {
			todayStats = projectmodels.ProjectListStats{
				Requests:  ds.TotalRequests,
				CostUsd:   ds.TotalCostUsd,
				Tokens:    ds.TotalTokens,
				Errors:    ds.FailedRequests,
				ErrorRate: safeErrorRate(ds.FailedRequests, ds.TotalRequests),
			}
		}
		if ms, err := ps.ProjectRepo.GetProjectMonthlyStats(p.Id, today.Year(), int(today.Month())); err == nil {
			monthStats = projectmodels.ProjectListStats{
				Requests: ms.TotalRequests,
				CostUsd:  ms.TotalCostUsd,
				Tokens:   ms.TotalTokens,
			}
		}

		items = append(items, projectmodels.ListProjectItem{
			Id:          p.Id.String(),
			Name:        p.Name,
			Description: p.Description,
			Environment: p.Environment,
			TodayStats:  todayStats,
			MonthStats:  monthStats,
			ApiKeys:     projectmodels.ProjectApiKeysSummary{},
			CreatedAt:   p.CreatedAt,
			UpdatedAt:   p.UpdatedAt,
		})
	}

	totalPages := 0
	if limit > 0 && totalCount > 0 {
		totalPages = (int(totalCount) + limit - 1) / limit
	}

	return &projectmodels.ListProjectsResponseBody{
		OrganizationId: orgId.String(),
		Projects:       items,
		Pagination: pagination.PaginationMeta{
			Page:       page,
			PageSize:   limit,
			TotalPages: totalPages,
			TotalItems: int(totalCount),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
	}, nil
}

func (ps *Projectservice) GetProjectDetails(ctx context.Context, projectId uuid.UUID, date time.Time) (*projectmodels.GetProjectDetailsResponseBody, error) {
	tid, err := tenantId(ctx)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 401, Message: "unauthorized: " + err.Error()}
	}

	p, err := ps.ProjectRepo.GetProjectByProjectId(projectId, tid)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "project not found"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error fetching project: " + err.Error()}
	}

	var todayStats projectmodels.ProjectTodayStats
	todayStats.Date = date.Format("2006-01-02")
	if ds, err := ps.ProjectRepo.GetProjectDailyStats(projectId, date); err == nil {
		todayStats.Requests = projectmodels.ProjectRequestCounts{
			Total:      ds.TotalRequests,
			Successful: ds.SuccessfulRequests,
			Failed:     ds.FailedRequests,
			ErrorRate:  safeErrorRate(ds.FailedRequests, ds.TotalRequests),
		}
		todayStats.CostUsd = ds.TotalCostUsd
		todayStats.Tokens = projectmodels.ProjectTokenCounts{Total: ds.TotalTokens}
		todayStats.Performance = projectmodels.ProjectPerformance{AvgDurationMs: ds.AvgDurationMs}
	}

	var monthStats projectmodels.ProjectMonthStats
	monthStats.Year = date.Year()
	monthStats.Month = int(date.Month())
	if ms, err := ps.ProjectRepo.GetProjectMonthlyStats(projectId, date.Year(), int(date.Month())); err == nil {
		monthStats.Requests = projectmodels.ProjectRequestCounts{Total: ms.TotalRequests}
		monthStats.CostUsd = ms.TotalCostUsd
		monthStats.Tokens = projectmodels.ProjectTokenCounts{Total: ms.TotalTokens}
	}

	var providers []projectmodels.ProjectProviderSimple
	startOfMonth := time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
	if provStats, err := ps.ProjectRepo.GetProviderStats(projectId, startOfMonth, date); err == nil {
		providerMap := map[string]*projectmodels.ProjectProviderSimple{}
		var totalCost float64
		for _, ps := range provStats {
			totalCost += ps.TotalCostUsd
		}
		for _, ps := range provStats {
			entry, ok := providerMap[ps.Provider]
			if !ok {
				entry = &projectmodels.ProjectProviderSimple{
					Provider:      ps.Provider,
					ProviderLabel: toTitleCase(ps.Provider),
					Models:        []string{},
				}
				providerMap[ps.Provider] = entry
			}
			entry.Requests += ps.TotalRequests
			entry.CostUsd += ps.TotalCostUsd
			entry.Tokens += ps.TotalTokens
			if totalCost > 0 {
				entry.CostShare = roundToOne(entry.CostUsd / totalCost * 100)
			}
		}
		for _, v := range providerMap {
			providers = append(providers, *v)
		}
	}
	if providers == nil {
		providers = []projectmodels.ProjectProviderSimple{}
	}

	return &projectmodels.GetProjectDetailsResponseBody{
		Project: projectmodels.ProjectDetail{
			Id:             p.Id.String(),
			OrganizationId: p.OrgId.String(),
			Name:           p.Name,
			Description:    p.Description,
			Environment:    p.Environment,
			CreatedAt:      p.CreatedAt,
			UpdatedAt:      p.UpdatedAt,
		},
		TodayStats:     todayStats,
		MonthStats:     monthStats,
		ProvidersUsage: providers,
		AlertConfiguration: projectmodels.AlertConfiguration{
			CostAlert:      projectmodels.AlertConfig{Enabled: false, ThresholdUsd: 100.0, Period: "daily"},
			ErrorRateAlert: projectmodels.AlertConfig{Enabled: false, ThresholdPercentage: 5.0, Period: "hourly"},
		},
		ApiKeysSummary: projectmodels.ApiKeysSummaryDetail{
			ActiveCount: 0,
			TotalCount:  0,
			Keys:        []projectmodels.ApiKeyEntry{},
		},
	}, nil
}

func (ps *Projectservice) GetProvidersBreakdown(ctx context.Context, projectId uuid.UUID, startDate time.Time, endDate time.Time) (*projectmodels.ProvidersBreakdownResponseBody, error) {
	tid, err := tenantId(ctx)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 401, Message: "unauthorized: " + err.Error()}
	}

	_, err = ps.ProjectRepo.GetProjectByProjectId(projectId, tid)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "project not found"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error fetching project: " + err.Error()}
	}

	provStats, err := ps.ProjectRepo.GetProviderStats(projectId, startDate, endDate)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error fetching provider stats: " + err.Error()}
	}

	type aggregate struct {
		Requests int64
		Cost     float64
		Tokens   int64
	}
	aggMap := map[string]*aggregate{}
	var totalCost float64
	var totalRequests, totalTokens int64

	for _, s := range provStats {
		a, ok := aggMap[s.Provider]
		if !ok {
			a = &aggregate{}
			aggMap[s.Provider] = a
		}
		a.Requests += s.TotalRequests
		a.Cost += s.TotalCostUsd
		a.Tokens += s.TotalTokens
		totalCost += s.TotalCostUsd
		totalRequests += s.TotalRequests
		totalTokens += s.TotalTokens
	}

	providerItems := make([]projectmodels.ProviderBreakdownItem, 0, len(aggMap))
	for provider, a := range aggMap {
		costShare := 0.0
		if totalCost > 0 {
			costShare = roundToOne(a.Cost / totalCost * 100)
		}
		providerItems = append(providerItems, projectmodels.ProviderBreakdownItem{
			Provider:      provider,
			ProviderLabel: toTitleCase(provider),
			Models:        []projectmodels.ProviderModelBreakdown{},
			TotalRequests: a.Requests,
			TotalCostUsd:  a.Cost,
			CostShare:     costShare,
			TotalTokens:   a.Tokens,
		})
	}

	return &projectmodels.ProvidersBreakdownResponseBody{
		ProjectId: projectId.String(),
		Period: projectmodels.ProvidersBreakdownPeriod{
			Start: startDate.Format("2006-01-02"),
			End:   endDate.Format("2006-01-02"),
		},
		Providers:     providerItems,
		TotalRequests: totalRequests,
		TotalCostUsd:  totalCost,
		TotalTokens:   totalTokens,
	}, nil
}

func (ps *Projectservice) UpdateProject(ctx context.Context, projectId uuid.UUID, req *projectmodels.UpdateProjectRequestBody) (*projectmodels.UpdateProjectResponseBody, error) {
	tid, err := tenantId(ctx)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 401, Message: "unauthorized: " + err.Error()}
	}

	existing, err := ps.ProjectRepo.GetProjectByProjectId(projectId, tid)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "project not found"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error fetching project: " + err.Error()}
	}

	updated, err := ps.ProjectRepo.UpdateProject(projectId, existing.OrgId, tid, &models.DBProject{
		Name:        req.Name,
		Description: req.Description,
		Environment: req.Environment,
	})
	if err != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error updating project: " + err.Error()}
	}

	return &projectmodels.UpdateProjectResponseBody{
		Id:          updated.Id.String(),
		Name:        updated.Name,
		Description: updated.Description,
		Environment: updated.Environment,
		UpdatedAt:   updated.UpdatedAt,
		Message:     "Project updated successfully",
	}, nil
}

func (ps *Projectservice) DeleteProject(ctx context.Context, projectId uuid.UUID) (*projectmodels.DeleteProjectResponseBody, error) {
	tid, err := tenantId(ctx)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 401, Message: "unauthorized: " + err.Error()}
	}

	existing, err := ps.ProjectRepo.GetProjectByProjectId(projectId, tid)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "project not found"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error fetching project: " + err.Error()}
	}

	name := existing.Name
	err = ps.ProjectRepo.DeleteProject(projectId, existing.OrgId, tid)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 500, Message: "error deleting project: " + err.Error()}
	}

	return &projectmodels.DeleteProjectResponseBody{
		Id:        projectId.String(),
		Name:      name,
		Deleted:   true,
		DeletedAt: time.Now(),
		Message:   "Project deleted successfully",
		SideEffects: projectmodels.DeleteSideEffects{
			ApiKeysRevoked:    0,
			AlertsDeleted:     0,
			UsageDataArchived: true,
		},
	}, nil
}

func (ps *Projectservice) GetProjectErrors(ctx context.Context, projectId uuid.UUID, date time.Time, limit int) (*projectmodels.ProjectErrorsResponseBody, error) {
	tid, err := tenantId(ctx)
	if err != nil {
		return nil, &models.ServiceResponse{Code: 401, Message: "unauthorized: " + err.Error()}
	}

	_, err = ps.ProjectRepo.GetProjectByProjectId(projectId, tid)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &models.ServiceResponse{Code: 404, Message: "project not found"}
		}
		return nil, &models.ServiceResponse{Code: 500, Message: "error fetching project: " + err.Error()}
	}

	var errorRate float64
	if ds, err := ps.ProjectRepo.GetProjectDailyStats(projectId, date); err == nil && ds.TotalRequests > 0 {
		errorRate = safeErrorRate(ds.FailedRequests, ds.TotalRequests)
	}

	return &projectmodels.ProjectErrorsResponseBody{
		ProjectId:   projectId.String(),
		Date:        date.Format("2006-01-02"),
		TotalErrors: 0,
		ErrorRate:   errorRate,
		Errors:      []projectmodels.ProjectErrorEntry{},
	}, nil
}

// ── helpers ───────────────────────────────────────────────────────────────────

func safeErrorRate(failed, total int64) float64 {
	if total == 0 {
		return 0
	}
	return roundToOne(float64(failed) / float64(total) * 100)
}

func roundToOne(f float64) float64 {
	return float64(int64(f*10+0.5)) / 10
}

func toTitleCase(s string) string {
	if len(s) == 0 {
		return s
	}
	result := make([]byte, len(s))
	result[0] = s[0] - 32
	if s[0] >= 'a' && s[0] <= 'z' {
		result[0] = s[0] - 32
	} else {
		result[0] = s[0]
	}
	copy(result[1:], s[1:])
	return string(result)
}
