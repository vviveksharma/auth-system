package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/internal/utils"
	dbmodels "github.com/vviveksharma/auth/models"
)

type AuthService interface {
	LoginUser(req *models.UserLoginRequest) (res *models.UserLoginResponse, err error)
}

type Auth struct {
	UserRepo  repo.UserRepositoryInterface
	LoginRepo repo.LoginRepositoryInterface
	RoleRepo  repo.RoleRepositoryInterface
}

func NewAuthService() (AuthService, error) {
	ser := &Auth{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (a *Auth) SetupRepo() error {
	var err error
	user, err := repo.NewUserRepository(db.DB)
	if err != nil {
		return err
	}
	a.UserRepo = user
	login, err := repo.NewLoginRepository(db.DB)
	if err != nil {
		return err
	}
	a.LoginRepo = login
	role, err := repo.NewRoleRepository(db.DB)
	if err != nil {
		return err
	}
	a.RoleRepo = role
	return nil
}

func (a *Auth) LoginUser(req *models.UserLoginRequest) (res *models.UserLoginResponse, err error) {
	userDetails, err := a.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "user with name doesnot exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while finding the userdetails: " + err.Error(),
			}
		}
	}
	if req.Password != userDetails.Password {
		return nil, &dbmodels.ServiceResponse{
			Code:    404,
			Message: "password doesn't exist or is expired",
		}
	}
	// Fetching the user roles
	if !contains(userDetails.Roles, req.Role) {
		return nil, &dbmodels.ServiceResponse{
			Code:    403,
			Message: "user does not have the required role",
		}
	}

	roleId, err := a.RoleRepo.FindRoleId(req.Role)
	if err != nil {
		if err.Error() == "record not found" || roleId == uuid.Nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "no role exist with this name",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while finding the roleDetails: " + err.Error(),
			}
		}
	}

	loginDetails, err := a.LoginRepo.GetUserById(userDetails.Id.String())
	if err != nil && err.Error() != "record not found" {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while finding the logindetails of the user: " + err.Error(),
		}
	}
	var tokenType string

	if loginDetails == nil {
		tokenType = "access"
	} else {
		tokenType = "refresh"
	}

	jwt, err := utils.CraeteJWT(userDetails.Id.String(), roleId.String(), tokenType)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a JWT: " + err.Error(),
		}
	}

	lerr := a.LoginRepo.Create(&dbmodels.DBLogin{
		UserId:    userDetails.Id,
		RoleId:    roleId,
		Token:     jwt,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Revoked:   false,
	})
	if lerr != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a login entry: " + lerr.Error(),
		}
	}

	return &models.UserLoginResponse{
		JWT: jwt,
	}, nil
}

func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
