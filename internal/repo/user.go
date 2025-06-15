package repo

import (
	"fmt"

	"github.com/google/uuid"
	dbmodels "github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.DBUser) error
	GetUserDetails(id uuid.UUID) (userDetails *models.DBUser, err error)
	GetUserByEmail(email string) (userDetails *models.DBUser, err error)
	UpdateUserFields(userID uuid.UUID, input *dbmodels.UpdateUserRequest) error
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

func (ur *UserRepository) GetUserDetails(id uuid.UUID) (userDetails *models.DBUser, err error) {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	user := transaction.First(&userDetails, models.DBUser{
		Id: id,
	})
	if user.Error != nil {
		return nil, user.Error
	}
	transaction.Commit()
	return userDetails, nil
}

func (ur *UserRepository) GetUserByEmail(email string) (userDetails *models.DBUser, err error) {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	user := transaction.First(&userDetails, models.DBUser{
		Email: email,
	})
	if user.Error != nil {
		return nil, user.Error
	}
	transaction.Commit()
	return userDetails, nil
}

func (r *UserRepository) UpdateUserFields(userID uuid.UUID, input *dbmodels.UpdateUserRequest) error {
	tx := r.DB.Begin()
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
