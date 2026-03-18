package orgrepo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type OrgRepositoryInterface interface {
	CreateOrg(req *models.DBOrganisation) error
	ListOrgTenant(tenantId uuid.UUID, page, pageSize int) (resp []*models.DBOrganisation, count int64, err error)
	GetOrgById(tenantId uuid.UUID, orgId uuid.UUID) (resp *models.DBOrganisation, err error)
	UpdateOrg(tenantId uuid.UUID, orgId uuid.UUID, req *models.DBOrganisation) (resp *models.DBOrganisation, err error)
	DeleteOrg(tenantId uuid.UUID, orgId uuid.UUID) error
	FindByConditons(req *models.DBOrganisation) (resp *models.DBOrganisation, err error)
}

type OrgRepository struct {
	DB *gorm.DB
}

func NewOrgRepository(db *gorm.DB) (OrgRepositoryInterface, error) {
	return &OrgRepository{db}, nil
}

func (org *OrgRepository) CreateOrg(req *models.DBOrganisation) error {
	transaction := org.DB.Begin()
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

func (org *OrgRepository) ListOrgTenant(tenantId uuid.UUID, page, pageSize int) (resp []*models.DBOrganisation, count int64, err error) {

	var totalCount int64
	var orgDetails []*models.DBOrganisation

	baseQuery := org.DB.Model(&models.DBOrganisation{}).Where("tenant_id = ?", tenantId)
	if err := baseQuery.Count(&totalCount).Error; err != nil {
		fmt.Printf("Error counting organisations in ListOrgTenant: %v\n", err)
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	if err := baseQuery.Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&orgDetails).Error; err != nil {
		fmt.Printf("Error fetching paginated organisations in ListOrgTenant: %v\n", err)
		return nil, 0, err
	}
	return orgDetails, totalCount, nil
}

func (org *OrgRepository) GetOrgById(tenantId uuid.UUID, orgId uuid.UUID) (resp *models.DBOrganisation, err error) {
	var orgDetails models.DBOrganisation
	query := org.DB.Model(&models.DBOrganisation{}).
		Where("tenant_id = ? AND id = ?", tenantId, orgId).
		Take(&orgDetails)
	if query.Error != nil {
		return nil, query.Error
	}
	return &orgDetails, nil
}

func (org *OrgRepository) UpdateOrg(tenantId uuid.UUID, orgId uuid.UUID, req *models.DBOrganisation) (resp *models.DBOrganisation, err error) {
	var orgDetails models.DBOrganisation
	query := org.DB.Model(&models.DBOrganisation{}).Where("tenant_id = ? AND id = ?", tenantId, orgId).
		Take(&orgDetails)
	if query.Error != nil {
		return nil, query.Error
	}

	updates := map[string]interface{}{}
	if req.Name != "" && orgDetails.Name != req.Name {
		updates["name"] = req.Name
	}
	if req.Slug != "" && orgDetails.Slug != req.Slug {
		updates["slug"] = req.Slug
	}
	if req.Description != "" && orgDetails.Description != req.Description {
		updates["description"] = req.Description
	}
	if req.IconUrl != "" && orgDetails.IconUrl != req.IconUrl {
		updates["icon_url"] = req.IconUrl
	}

	if len(updates) == 0 {
		return &orgDetails, nil
	}

	updatedAt := time.Now()
	updates["updated_at"] = updatedAt

	updateQuery := org.DB.Model(&models.DBOrganisation{}).
		Where("tenant_id = ? AND id = ?", tenantId, orgId).
		Updates(updates)
	if updateQuery.Error != nil {
		return nil, updateQuery.Error
	}

	if name, ok := updates["name"].(string); ok {
		orgDetails.Name = name
	}
	if slug, ok := updates["slug"].(string); ok {
		orgDetails.Slug = slug
	}
	if description, ok := updates["description"].(string); ok {
		orgDetails.Description = description
	}
	if iconUrl, ok := updates["icon_url"].(string); ok {
		orgDetails.IconUrl = iconUrl
	}
	orgDetails.UpdatedAt = updatedAt

	return &orgDetails, nil
}

func (org *OrgRepository) DeleteOrg(tenantId uuid.UUID, orgId uuid.UUID) error {
	query := org.DB.Unscoped().Where("tenant_id = ? AND id = ?", tenantId, orgId).Delete(&models.DBOrganisation{})
	if query.Error != nil {
		return query.Error
	}
	return nil
}

func (org *OrgRepository) FindByConditons(req *models.DBOrganisation) (resp *models.DBOrganisation, err error) {
	var orgDetails models.DBOrganisation
	if req == nil {
		return nil, gorm.ErrInvalidData
	}

	conditions := map[string]interface{}{}
	if req.Id != uuid.Nil {
		conditions["id"] = req.Id
	}
	if req.TenantId != uuid.Nil {
		conditions["tenant_id"] = req.TenantId
	}
	if req.Name != "" {
		conditions["name"] = req.Name
	}
	if req.Slug != "" {
		conditions["slug"] = req.Slug
	}
	if req.Status != "" {
		conditions["status"] = req.Status
	}

	if len(conditions) == 0 {
		return nil, gorm.ErrMissingWhereClause
	}

	query := org.DB.Model(&models.DBOrganisation{}).Where(conditions).Take(&orgDetails)
	if query.Error != nil {
		return nil, query.Error
	}
	return &orgDetails, nil
}
