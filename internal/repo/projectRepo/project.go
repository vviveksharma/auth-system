package projectrepo

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type ProjectRepositryInterface interface {
	Create(req *models.DBProject) error
	ListOrgProject(orgId uuid.UUID, tenantId uuid.UUID, page int, pageSize int) ([]*models.DBProject, int64, error)
	GetProjectById(projectId uuid.UUID, orgId uuid.UUID, tenantId uuid.UUID) (*models.DBProject, error)
	GetProjectByProjectId(projectId uuid.UUID, tenantId uuid.UUID) (*models.DBProject, error)
	DeleteProject(projectId uuid.UUID, orgId uuid.UUID, tenantId uuid.UUID) error
	GetProjectByName(orgId uuid.UUID, tenantId uuid.UUID, projectName string) (*models.DBProject, error)
	UpdateProject(projectId uuid.UUID, orgId uuid.UUID, tenantId uuid.UUID, fields *models.DBProject) (*models.DBProject, error)
	GetProjectDailyStats(projectId uuid.UUID, date time.Time) (*models.DBProjectDailyStats, error)
	GetProjectMonthlyStats(projectId uuid.UUID, year int, month int) (*models.DBProjectMonthlyStats, error)
	GetProviderStats(projectId uuid.UUID, startDate time.Time, endDate time.Time) ([]*models.DBProviderDailyStats, error)
}

type ProjectRepositry struct {
	DB *gorm.DB
}

func NewProjectReposistry(db *gorm.DB) (ProjectRepositryInterface, error) {
	return &ProjectRepositry{
		DB: db,
	}, nil
}

func (p *ProjectRepositry) Create(req *models.DBProject) error {
	transaction := p.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	create := transaction.Create(&req)
	if create.Error != nil {
		return create.Error
	}
	transaction.Commit()
	return nil
}

func (p *ProjectRepositry) ListOrgProject(orgId uuid.UUID, tenantId uuid.UUID, page int, pageSize int) ([]*models.DBProject, int64, error) {
	var totalCount int64
	var response []*models.DBProject

	transaction := p.DB.Begin()
	if transaction.Error != nil {
		return nil, 0, transaction.Error
	}
	defer transaction.Rollback()

	baseQuery := transaction.Model(&models.DBProject{}).Where("org_id = ? AND tenant_id = ?", orgId, tenantId)

	if err := baseQuery.Count(&totalCount).Error; err != nil {
		log.Printf("Error counting org projects: %v", err)
		return nil, 0, fmt.Errorf("error counting projects for org: %v", err)
	}

	if totalCount == 0 {
		if err := transaction.Commit().Error; err != nil {
			return nil, 0, err
		}
		return []*models.DBProject{}, 0, nil
	}

	offset := (page - 1) * pageSize

	if err := transaction.Model(&models.DBProject{}).
		Where("org_id = ? AND tenant_id = ?", orgId, tenantId).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&response).Error; err != nil {
		log.Printf("Error fetching org projects: %v", err)
		return nil, 0, fmt.Errorf("error fetching projects: %v", err)
	}

	if err := transaction.Commit().Error; err != nil {
		return nil, 0, err
	}

	return response, totalCount, nil
}

func (p *ProjectRepositry) GetProjectById(projectId uuid.UUID, orgId uuid.UUID, tenantId uuid.UUID) (*models.DBProject, error) {
	transaction := p.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var response models.DBProject
	getQuery := transaction.Where(&models.DBProject{}).Where("id = ? AND tenant_id = ? AND org_id = ?", projectId, tenantId, orgId).Take(&response)
	if getQuery.Error != nil {
		return nil, getQuery.Error
	}
	transaction.Commit()
	return &response, nil
}

func (p *ProjectRepositry) DeleteProject(projectId uuid.UUID, orgId uuid.UUID, tenantId uuid.UUID) error {
	transaction := p.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	deleteQuery := transaction.Where(&models.DBProject{}).Unscoped().Where("id = ? AND tenant_id = ? AND org_id = ?", projectId, tenantId, orgId).Delete(&models.DBProject{})
	if deleteQuery.Error != nil {
		return deleteQuery.Error
	}
	return transaction.Commit().Error
}

func (p *ProjectRepositry) GetProjectByName(orgId uuid.UUID, tenantId uuid.UUID, projectName string) (*models.DBProject, error) {
	transaction := p.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	var response *models.DBProject
	getByNameQuery := transaction.Where(&models.DBProject{}).Where("tenant_id = ? AND org_id = ? AND name = ?", tenantId, orgId, projectName).First(&response)
	if getByNameQuery.Error != nil {
		return nil, getByNameQuery.Error
	}
	return response, nil
}

func (p *ProjectRepositry) GetProjectByProjectId(projectId uuid.UUID, tenantId uuid.UUID) (*models.DBProject, error) {
	var response models.DBProject
	err := p.DB.Where("id = ? AND tenant_id = ?", projectId, tenantId).Take(&response).Error
	if err != nil {
		return nil, err
	}
	return &response, nil
}

func (p *ProjectRepositry) UpdateProject(projectId uuid.UUID, orgId uuid.UUID, tenantId uuid.UUID, fields *models.DBProject) (*models.DBProject, error) {
	var result models.DBProject
	now := time.Now()
	fields.UpdatedAt = &now
	err := p.DB.Model(&models.DBProject{}).
		Where("id = ? AND org_id = ? AND tenant_id = ?", projectId, orgId, tenantId).
		Updates(map[string]interface{}{
			"name":        fields.Name,
			"description": fields.Description,
			"environment": fields.Environment,
			"updated_at":  now,
		}).Error
	if err != nil {
		return nil, err
	}
	if err := p.DB.Where("id = ? AND org_id = ? AND tenant_id = ?", projectId, orgId, tenantId).First(&result).Error; err != nil {
		return nil, err
	}
	return &result, nil
}

func (p *ProjectRepositry) GetProjectDailyStats(projectId uuid.UUID, date time.Time) (*models.DBProjectDailyStats, error) {
	var stats models.DBProjectDailyStats
	err := p.DB.Where("project_id = ? AND date = ?", projectId, date.Format("2006-01-02")).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (p *ProjectRepositry) GetProjectMonthlyStats(projectId uuid.UUID, year int, month int) (*models.DBProjectMonthlyStats, error) {
	var stats models.DBProjectMonthlyStats
	err := p.DB.Where("project_id = ? AND year = ? AND month = ?", projectId, year, month).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, nil
}

func (p *ProjectRepositry) GetProviderStats(projectId uuid.UUID, startDate time.Time, endDate time.Time) ([]*models.DBProviderDailyStats, error) {
	var stats []*models.DBProviderDailyStats
	err := p.DB.Where("project_id = ? AND date BETWEEN ? AND ?", projectId, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Find(&stats).Error
	if err != nil {
		return nil, err
	}
	return stats, nil
}

