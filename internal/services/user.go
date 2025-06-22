package services

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
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
	err = u.UserRepo.CreateUser(&dbmodels.DBUser{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Roles:    []string{"guest"},
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
