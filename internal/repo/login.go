package repo

import (
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type LoginRepositoryInterface interface {
	Create(req *models.DBLogin) error
	GetUserById(id string) (loginDetails *models.DBLogin, err error)
}

type LoginRepository struct {
	DB *gorm.DB
}

func NewLoginRepository(db *gorm.DB) (LoginRepositoryInterface, error) {
	return &LoginRepository{DB: db}, nil
}

func (r *LoginRepository) Create(req *models.DBLogin) error {
	transaction := r.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	newUser := transaction.Create(&req)
	if newUser.Error != nil {
		return newUser.Error
	}
	transaction.Commit()
	return nil
}

func (ur *LoginRepository) GetUserById(id string) (loginDetails *models.DBLogin, err error) {
	transaction := ur.DB.Begin()
	if transaction.Error != nil {
		return nil, transaction.Error
	}
	defer transaction.Rollback()
	user := transaction.First(&loginDetails, &models.DBLogin{
		UserId: uuid.MustParse(id),
	})
	if user.Error != nil {
		return nil, user.Error
	}
	transaction.Commit()
	return loginDetails, nil
}
