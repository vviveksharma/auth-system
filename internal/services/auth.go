package services

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/internal/utils"
	dbmodels "github.com/vviveksharma/auth/models"
)

type AuthService interface {
	LoginUser(req *models.UserLoginRequest, ctx context.Context) (res *models.UserLoginResponse, err error)
	RefreshToken(id string, roleId string, ctx context.Context) (res *models.UserLoginResponse, err error)
	LogoutUser(userId uuid.UUID, ctx context.Context) (*models.LogoutUserResponse, error)
}

type Auth struct {
	UserRepo       repo.UserRepositoryInterface
	LoginRepo      repo.LoginRepositoryInterface
	RoleRepo       repo.RoleRepositoryInterface
	TokenRepo      repo.TokenRepositoryInterface
	ResetTokenRepo repo.ResetTokenRepositoryInterface
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
	rToken, err := repo.NewResetTokenRepository(db.DB)
	if err != nil {
		return err
	}
	a.ResetTokenRepo = rToken
	Token, err := repo.NewTokenRepository(db.DB)
	if err != nil {
		return err
	}
	a.TokenRepo = Token
	return nil
}

func (a *Auth) LoginUser(req *models.UserLoginRequest, ctx context.Context) (res *models.UserLoginResponse, err error) {
	tenant_id := ctx.Value("tenant_id").(string)
	fmt.Println("the tenant id : ", tenant_id)
	userDetails, err := a.UserRepo.GetUserByEmail(req.Email, uuid.MustParse(tenant_id))
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
	flag, err := utils.ComparePassword(req.Password, userDetails.Password, userDetails.Salt, utils.DefaultParams)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    401,
			Message: "error while comparing password: " + err.Error(),
		}
	}
	if !flag {
		return nil, &dbmodels.ServiceResponse{
			Code:    401,
			Message: "Invalid password",
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

	jwt, err := utils.CreateJWT(userDetails.Id.String(), roleId.String(), userDetails.TenantId.String(), tokenType)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a JWT: " + err.Error(),
		}
	}

	if tokenType == "refresh" {
		err := a.LoginRepo.UpdateUserToken(loginDetails.Id.String(), jwt)
		if err != nil {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while creating a JWT: " + err.Error(),
			}
		}
		return &models.UserLoginResponse{
			JWT: jwt,
		}, nil
	}

	lerr := a.LoginRepo.Create(&dbmodels.DBLogin{
		UserId:    userDetails.Id,
		RoleId:    roleId,
		RoleName:  req.Role,
		JWTToken:  jwt,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(30 * time.Minute),
		Revoked:   false,
		TenantId:  uuid.MustParse(tenant_id),
	})
	if lerr != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a login entry: " + lerr.Error(),
		}
	}
	// Check if while login on the guest role is there update that with user
	return &models.UserLoginResponse{
		JWT: jwt,
	}, nil
}

func (a *Auth) RefreshToken(id string, roleId string, ctx context.Context) (res *models.UserLoginResponse, err error) {
	loginDetails, err := a.LoginRepo.GetUserById(id)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while fetching then login entry: " + err.Error(),
		}
	}
	jwt, err := utils.CreateJWT(loginDetails.Id.String(), roleId, loginDetails.TenantId.String(), "refresh")
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a JWT: " + err.Error(),
		}
	}
	err = a.LoginRepo.UpdateUserToken(loginDetails.Id.String(), jwt)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while refreshing the token: " + err.Error(),
		}
	}
	return &models.UserLoginResponse{
		JWT: jwt,
	}, nil
}

func (a *Auth) LogoutUser(userId uuid.UUID, ctx context.Context) (*models.LogoutUserResponse, error) {
	err := a.LoginRepo.Logout(userId)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to log out user with id %s: %v", userId.String(), err),
		}
	}
	return &models.LogoutUserResponse{
		Message: fmt.Sprintf("User with ID %s has been logged out successfully.", userId.String()),
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
