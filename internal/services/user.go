package services

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	"github.com/vviveksharma/auth/internal/utils"
	dbmodels "github.com/vviveksharma/auth/models"
)

type UserService interface {
	CreateUser(req *models.UserRequest) (*models.UserResponse, error)
	GetUserDetails(req *models.GetUserDetailsRequest) (*models.UserDetailsResponse, error)
	UpdateUserDetails(req *models.UpdateUserRequest, userId string) (*models.UpdateUserResponse, error)
	GetUserById(userId string) (*models.GetUserByIdResponse, error)
	AssignUserRole(req *models.AssignRoleRequest, userId string) (*models.AssignRoleResponse, error)
}

type User struct {
	UserRepo       repo.UserRepositoryInterface
	ResetTokenRepo repo.ResetTokenRepositoryInterface
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

	token, err := repo.NewResetTokenRepository(db.DB)
	if err != nil {
		return err
	}
	u.ResetTokenRepo = token
	return nil
}

func (u *User) CreateUser(req *models.UserRequest) (*models.UserResponse, error) {
	userDetails, err := u.UserRepo.GetUserByEmail(req.Email)
	if err != nil && err.Error() != "record not found" {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating the userdetails: " + err.Error(),
		}
	}
	if userDetails != nil && userDetails.Email == req.Email {
		return nil, &dbmodels.ServiceResponse{
			Code:    404,
			Message: "user with name exist please login",
		}
	}
	// Storing the hash password
	pass, salt, err := utils.GeneratePasswordHash(req.Password, utils.DefaultParams)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    423,
			Message: "error while generating and storing the password: " + err.Error(),
		}
	}
	id := uuid.MustParse("dae760ab-0a7f-4cbd-8603-def85ad8e430")
	err = u.UserRepo.CreateUser(&dbmodels.DBUser{
		TenantId: id,
		Name:     req.Name,
		Email:    req.Email,
		Password: pass,
		Salt:     salt,
		Roles:    []string{"admin"},
	})
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating the userdetails: " + err.Error(),
		}
	}
	return &models.UserResponse{
		Message: "user created succesfuly with email " + req.Email,
	}, nil
}

func (u *User) GetUserDetails(req *models.GetUserDetailsRequest) (*models.UserDetailsResponse, error) {
	userDetails, err := u.UserRepo.GetUserDetails(uuid.MustParse(req.Id))
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
	return &models.UserDetailsResponse{
		Email: userDetails.Email,
		Name:  userDetails.Name,
		Role:  userDetails.Roles,
	}, nil
}

func (u *User) UpdateUserDetails(req *models.UpdateUserRequest, userId string) (*models.UpdateUserResponse, error) {
	userDetails, err := u.UserRepo.GetUserDetails(uuid.MustParse(userId))
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
	fmt.Println("the user present with the details: ", userDetails)
	fmt.Println("the request : ", req.Email, req.Name, req.Password)
	err = u.UserRepo.UpdateUserFields(uuid.MustParse(userId), req)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the user details: " + err.Error(),
		}
	}
	return &models.UpdateUserResponse{
		Message: "user updated successfully",
	}, nil
}

func (u *User) GetUserById(userId string) (*models.GetUserByIdResponse, error) {
	userDetails, err := u.UserRepo.GetUserDetails(uuid.MustParse(userId))
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
	return &models.GetUserByIdResponse{
		Name:  userDetails.Name,
		Email: userDetails.Email,
		Role:  userDetails.Roles,
	}, nil
}

func (u *User) AssignUserRole(req *models.AssignRoleRequest, userId string) (*models.AssignRoleResponse, error) {
	userDetails, err := u.UserRepo.GetUserDetails(uuid.MustParse(userId))
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
	err = u.UserRepo.UpdateUserRoles(userDetails.Id, req.Role)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while assigning the role: " + err.Error(),
		}
	}
	return &models.AssignRoleResponse{
		Message: "User role updated successfully",
	}, nil
}

func (u *User) ResetPassword(req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error) {
	userDetails, err := u.UserRepo.GetUserByEmail(req.Email)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    400,
				Message: "Unable to find the user with the provided email: " + req.Email,
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("internal error while searching for user with email '%s': %s", req.Email, err.Error()),
			}
		}
	}
	// Create a unique token valid for 5 minutes
	tokenErr := u.ResetTokenRepo.Create(&dbmodels.DBResetToken{
		UserId:    userDetails.Id,
		TenantId:  userDetails.TenantId,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		IsActive:  true,
	})
	if tokenErr != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to generate password reset token for user '%s': %s", req.Email, tokenErr.Error()),
		}
	}
	return &models.ResetPasswordResponse{}, nil
}
