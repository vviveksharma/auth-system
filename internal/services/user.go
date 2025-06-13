package services

import (
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type UserService interface {
	CreateUser(req *models.UserRequest) (*models.UserResponse, error)
}

type User struct {
	UserRepo repo.UserRepositoryInterface
}

func NewUserService() (UserService, error) {
	ser := &User{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (u *User) SetupRepo() error {
	var err error
	user, err := repo.NewUserRepository(db.DB)
	if err != nil {
		return err
	}
	u.UserRepo = user
	return nil
}

func (u *User) CreateUser(req *models.UserRequest) (*models.UserResponse, error) {
	err := u.UserRepo.CreateUser(&dbmodels.DBUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Role:     "",
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "user with name doesnot exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while creating the userdetails: " + err.Error(),
			}
		}
	}
	return &models.UserResponse{
		Message: "user created succesfuly with email " + req.Email,
	}, nil
}
