package repo

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *models.DBUser) error
	GetUserDetails(id uuid.UUID) (userDetails *models.DBUser, err error)
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
