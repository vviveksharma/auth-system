package repo

import (
	"fmt"
	"log"

	"github.com/google/uuid"
	dbmodels "github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.DBUser) error
	GetUserDetails(conditions models.DBUser) (userDetails *models.DBUser, err error)
	GetUserByEmail(email string, tenantId uuid.UUID) (userDetails *models.DBUser, err error)
	UpdateUserFields(userID uuid.UUID, input *dbmodels.UpdateUserRequest) error
	UpdateUserRoles(userId uuid.UUID, role string) error
	UpdatePassword(userId uuid.UUID, password string) error
	ListUsersPaginated(page int, pageSize int, tenantId uuid.UUID, status string) ([]*models.DBUser, int64, error)
	ChangeStatus(flag bool, id uuid.UUID) error
}

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) (UserRepositoryInterface, error) {
	return &UserRepository{DB: db}, nil
}

func (ur *UserRepository) CreateUser(user *models.DBUser) error {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	fmt.Println("the dal-layer: ", user)
	newUser := transaction.Create(&user)
	if newUser.Error != nil {
		return newUser.Error
	}
	fmt.Println("the error:", newUser.Error)
	transaction.Commit()
	return nil
}

func (ur *UserRepository) GetUserDetails(conditions models.DBUser) (userDetails *models.DBUser, err error) {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	user := transaction.First(&userDetails, &conditions)
	if user.Error != nil {
		return nil, user.Error
	}
	transaction.Commit()
	return userDetails, nil
}

func (ur *UserRepository) GetUserByEmail(email string, tenantId uuid.UUID) (userDetails *models.DBUser, err error) {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	user := transaction.First(&userDetails, models.DBUser{
		Email:    email,
		TenantId: tenantId,
	})
	if user.Error != nil {
		return nil, user.Error
	}
	transaction.Commit()
	return userDetails, nil
}

func (ur *UserRepository) UpdateUserFields(userID uuid.UUID, input *dbmodels.UpdateUserRequest) error {
	tx := ur.DB.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	updates := map[string]interface{}{}

	if input.Name != nil {
		updates["name"] = *input.Name
	}
	if input.Email != nil {
		updates["email"] = *input.Email
	}
	if input.Password != nil {
		updates["password"] = *input.Password
	}

	if len(updates) == 0 {
		tx.Rollback()
		return nil
	}

	fmt.Println("the update models is:", updates)

	if err := tx.Model(&models.DBUser{}).
		Where("id = ?", userID).
		Updates(updates).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (ur *UserRepository) UpdateUserRoles(userId uuid.UUID, role string) error {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	userDetails := models.DBUser{}
	user := transaction.First(&userDetails, models.DBUser{
		Id: userId,
	})
	if user.Error != nil {
		return user.Error
	}
	log.Println("found out the user")
	userDetails.Roles = append(userDetails.Roles, role)
	update := transaction.Save(&userDetails)
	if update.Error != nil {
		return update.Error
	}
	transaction.Commit()
	return nil
}

func (ur *UserRepository) UpdatePassword(userId uuid.UUID, password string) error {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	update := transaction.Model(models.DBUser{}).Where("id = ?", userId).Updates(map[string]interface{}{
		"password": password,
	})
	if update.Error != nil {
		return update.Error
	}
	return nil
}

func (tu *UserRepository) ListUsersPaginated(page int, pageSize int, tenantId uuid.UUID, status string) ([]*models.DBUser, int64, error) {
	log.Printf("🔍 [REPO] ListUsersPaginated called with: page=%d, pageSize=%d, tenantId=%s, status=%s",
		page, pageSize, tenantId.String(), status)

	var totalCount int64
	var users []*models.DBUser

	var is_active bool
	if status == "enabled" {
		is_active = true
	} else {
		is_active = false
	}
	log.Printf("🔍 [REPO] Converted status '%s' to is_active=%v", status, is_active)

	// Count total matching users
	countQuery := tu.DB.Model(&models.DBUser{}).
		Where("tenant_id = ?", tenantId).
		Where("status = ?", is_active)

	log.Printf("🔍 [REPO] Count query SQL: %v", countQuery.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Count(&totalCount)
	}))

	if err := countQuery.Count(&totalCount).Error; err != nil {
		log.Printf("❌ [REPO] Error counting the users: %v", err)
		return nil, 0, fmt.Errorf("error counting the users present for this tenant: %v", err)
	}
	log.Printf("✅ [REPO] Total users found: %d", totalCount)

	offset := (page - 1) * pageSize
	log.Printf("🔍 [REPO] Calculated offset: %d (page=%d, pageSize=%d)", offset, page, pageSize)

	// Fetch paginated users
	fetchQuery := tu.DB.Model(&models.DBUser{}).
		Where("tenant_id = ?", tenantId).
		Where("status = ?", is_active).
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset)

	log.Printf("🔍 [REPO] Fetch query SQL: %v", fetchQuery.ToSQL(func(tx *gorm.DB) *gorm.DB {
		return tx.Find(&users)
	}))

	if err := fetchQuery.Find(&users).Error; err != nil {
		log.Printf("❌ [REPO] Error fetching paginated users for the tenant: %v", err)
		return nil, 0, fmt.Errorf("error fetching users for this tenant: %w", err)
	}

	log.Printf("✅ [REPO] Successfully fetched %d users for this tenant", len(users))
	for i, user := range users {
		log.Printf("   [REPO] User %d: ID=%s, Email=%s, Name=%s, Status=%v",
			i+1, user.Id.String(), user.Email, user.Name, user.Status)
	}

	return users, totalCount, nil
}

func (ur *UserRepository) DeleteUser(id uuid.UUID) error {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	usersDelete := transaction.Model(models.DBUser{}).Where("id = ?", id).Delete(models.DBUser{
		Id: id,
	})
	if usersDelete.Error != nil {
		return usersDelete.Error
	}
	return nil
}

func (ur *UserRepository) ChangeStatus(flag bool, id uuid.UUID) error {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	if !flag { // disable the role
		update := transaction.Model(&models.DBUser{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status": false,
		})
		if update.Error != nil {
			return update.Error
		}
	} else { // enable the role
		update := transaction.Model(&models.DBUser{}).Where("id = ?", id).Updates(map[string]interface{}{
			"status": true,
		})
		if update.Error != nil {
			return update.Error
		}
	}
	return nil
}
