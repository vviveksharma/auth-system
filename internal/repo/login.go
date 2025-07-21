package repo

import (
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/models"
	"gorm.io/gorm"
)

type LoginRepositoryInterface interface {
	Create(req *models.DBLogin) error
	GetUserById(id string) (loginDetails *models.DBLogin, err error)
	UpdateUserToken(id string, jwt string) error
	DeleteToken(id string) error
}

type LoginRepository struct {
	DB *gorm.DB
}

func NewLoginRepository(db *gorm.DB) (LoginRepositoryInterface, error) {
	return &LoginRepository{DB: db}, nil
}

func (l *LoginRepository) Create(req *models.DBLogin) error {
	transaction := l.DB.Begin()
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

func (l *LoginRepository) GetUserById(id string) (loginDetails *models.DBLogin, err error) {
	transaction := l.DB.Begin()
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

func (l *LoginRepository) UpdateUserToken(id string, jwt string) error {
	transaction := l.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	var loginDetails *models.DBLogin
	login := transaction.Where("id = ?", uuid.MustParse(id)).First(&loginDetails)
	if login.Error != nil {
		return login.Error
	}
	if err := transaction.Model(&models.DBLogin{}).Where("id = ?", uuid.MustParse(id)).Updates(map[string]interface{}{
		"token":      jwt,
		"issued_at":  time.Now(),
		"expires_at": time.Now().Add(30 * time.Minute),
	}).Error; err != nil {
		return err
	}
	transaction.Commit()
	return nil
}

func (l *LoginRepository) DeleteToken(id string) error {
	transaction := l.DB.Begin()
	if transaction.Error != nil {
		return transaction.Error
	}
	defer transaction.Rollback()
	if err := transaction.Model(&models.DBLogin{}).Where("id = ? ", uuid.MustParse(id)).Updates(map[string]interface{}{
		"revoked":    true,
		"expires_at": time.Now(),
	}).Error; err != nil {
		return err
	}
	transaction.Commit()
	return nil
}
