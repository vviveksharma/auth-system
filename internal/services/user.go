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

type UserService interface {
	GetUserDetails(ctx context.Context, req *models.GetUserDetailsRequest) (*models.UserDetailsResponse, error)
	UpdateUserDetails(ctx context.Context, req *models.UpdateUserRequest, userId string) (*models.UpdateUserResponse, error)
	GetUserById(ctx context.Context, userId string) (*models.GetUserByIdResponse, error)
	AssignUserRole(ctx context.Context, req *models.AssignRoleRequest, userId string) (*models.AssignRoleResponse, error)
	RegisterUser(ctx context.Context, req *models.UserRequest) (*models.UserResponse, error)
	ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error)
	SetPassword(ctx context.Context, req *models.UserVerifyOTPRequest) (*models.UserVerifyOTPRequest, error)
	DeleteUser(ctx context.Context, userId uuid.UUID) (*models.DeleteUserResponse, error)
	ListUsers(ctx context.Context) (resp []*models.ListUsersResponse, err error)
	EnableUser(ctx context.Context, userId uuid.UUID) (*models.EnableUserResponse, error)
	DisbaleUser(ctx context.Context, userId uuid.UUID) (*models.DisableUserResponse, error)
	GetUserRole(ctx context.Context, userId uuid.UUID) (*models.GetRoleDetailsUser, error)
}

type User struct {
	UserRepo       repo.UserRepositoryInterface
	TokenRepo      repo.TokenRepository
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

	rtoken, err := repo.NewResetTokenRepository(db.DB)
	if err != nil {
		return err
	}
	u.ResetTokenRepo = rtoken
	return nil
}

func (u *User) RegisterUser(ctx context.Context, req *models.UserRequest) (*models.UserResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserByEmail(req.Email, uuid.MustParse(tenantId))
	if err != nil {
		if err.Error() != "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the user details: " + err.Error(),
			}
		}
	}
	if userDetails != nil && userDetails.Email == req.Email {
		return nil, &dbmodels.ServiceResponse{
			Code:    423,
			Message: "this user already exist proceed to login",
		}
	}
	hashedpassword, salt, err := utils.GeneratePasswordHash(req.Password, utils.DefaultParams)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while generating the password for the user",
		}
	}
	err = u.UserRepo.CreateUser(&dbmodels.DBUser{
		TenantId:  uuid.MustParse(tenantId),
		CreatedAt: time.Now(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedpassword,
		Salt:      salt,
		Roles:     []string{"guest"},
	})
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while registering the user: " + err.Error(),
		}
	}
	return &models.UserResponse{
		Message: "user registered successfully",
	}, nil
}

func (u *User) GetUserDetails(ctx context.Context, req *models.GetUserDetailsRequest) (*models.UserDetailsResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       uuid.MustParse(req.Id),
		TenantId: uuid.MustParse(tenantId),
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
	return &models.UserDetailsResponse{
		Email: userDetails.Email,
		Name:  userDetails.Name,
		Role:  userDetails.Roles,
	}, nil
}

func (u *User) UpdateUserDetails(ctx context.Context, req *models.UpdateUserRequest, userId string) (*models.UpdateUserResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       uuid.MustParse(userId),
		TenantId: uuid.MustParse(tenantId),
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

func (u *User) GetUserById(ctx context.Context, userId string) (*models.GetUserByIdResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       uuid.MustParse(userId),
		TenantId: uuid.MustParse(tenantId),
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
	return &models.GetUserByIdResponse{
		Name:  userDetails.Name,
		Email: userDetails.Email,
		Role:  userDetails.Roles,
	}, nil
}

func (u *User) AssignUserRole(ctx context.Context, req *models.AssignRoleRequest, userId string) (*models.AssignRoleResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       uuid.MustParse(userId),
		TenantId: uuid.MustParse(tenantId),
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

func (u *User) ResetPassword(ctx context.Context, req *models.ResetPasswordRequest) (*models.ResetPasswordResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserByEmail(req.Email, uuid.MustParse(tenantId))
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
	token, tokenErr := u.ResetTokenRepo.Create(&dbmodels.DBResetToken{
		UserId:    userDetails.Id,
		TenantId:  userDetails.TenantId,
		ExpiresAt: time.Now().Add(15 * time.Minute),
		IsActive:  true,
	})
	if tokenErr != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to generate password reset token for user '%s': %s", req.Email, tokenErr.Error()),
		}
	}
	return &models.ResetPasswordResponse{
		Message: "otp for the user: " + token.String(),
	}, nil
}

func (u *User) SetPassword(ctx context.Context, req *models.UserVerifyOTPRequest) (*models.UserVerifyOTPRequest, error) {
	tenantId := ctx.Value("tenant_id").(string)
	istoken, err := u.ResetTokenRepo.VerifyOTP(req.OTP)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while verifying the token: " + err.Error(),
		}
	}
	if !istoken {
		return nil, &dbmodels.ServiceResponse{
			Code:    423,
			Message: "token is already expired please try again",
		}
	}
	// lets set new pass
	userDetails, err := u.UserRepo.GetUserByEmail(req.Email, uuid.MustParse(tenantId))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "user with name doesnot exist",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: "error while fetching the userdetails: " + err.Error(),
			}
		}
	}
	hashedpassword, err := utils.GeneratePassword(req.NewPassword, utils.DefaultParams, userDetails.Salt)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while generating the hash for the user: " + err.Error(),
		}
	}
	err = u.UserRepo.UpdatePassword(userDetails.Id, hashedpassword)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while updating the password for the user: " + err.Error(),
		}
	}
	return nil, nil
}

func (u *User) DeleteUser(ctx context.Context, userId uuid.UUID) (*models.DeleteUserResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       userId,
		TenantId: uuid.MustParse(tenantId),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: fmt.Sprintf("No user found with the provided ID: %s.", userId.String()),
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An error occurred while retrieving user details for ID %s: %s.", userId.String(), err.Error()),
			}
		}
	}
	err = u.UserRepo.DeleteUser(userDetails.Id)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to delete user with ID %s: %s.", userId.String(), err.Error()),
		}
	}
	return &models.DeleteUserResponse{
		Message: fmt.Sprintf("User with ID %s has been deleted successfully.", userId.String()),
	}, nil
}

func (u *User) ListUsers(ctx context.Context) (resp []*models.ListUsersResponse, err error) {
	tenant_id := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.ListUsers(uuid.MustParse(tenant_id))
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: fmt.Sprintf("No user found with the provided tenanatID: %s.", tenant_id),
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("An error occurred while retrieving user details for given tenant ID %s: %s.", tenant_id, err.Error()),
			}
		}
	}
	for _, user := range userDetails {
		userdetails := &models.ListUsersResponse{
			Id:        user.Id,
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			Roles:     user.Roles,
		}
		resp = append(resp, userdetails)
	}
	return resp, nil
}

func (u *User) EnableUser(ctx context.Context, userId uuid.UUID) (*models.EnableUserResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       userId,
		TenantId: uuid.MustParse(tenantId),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: fmt.Sprintf("User not found for tenant %s and user ID %s", tenantId, userId.String()),
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("Failed to retrieve user details for tenant %s and user ID %s: %s", tenantId, userId.String(), err.Error()),
			}
		}
	}
	if userDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    400,
			Message: fmt.Sprintf("User %s is already enabled for tenant %s", userId.String(), tenantId),
		}
	}
	// Enable the disabled user
	err = u.UserRepo.ChangeStatus(true, userDetails.Id)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to enable user %s for tenant %s: %s", userId.String(), tenantId, err.Error()),
		}
	}
	return &models.EnableUserResponse{
		Message: fmt.Sprintf("User %s successfully enabled for tenant %s", userId.String(), tenantId),
	}, nil
}

func (u *User) DisbaleUser(ctx context.Context, userId uuid.UUID) (*models.DisableUserResponse, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		Id:       userId,
		TenantId: uuid.MustParse(tenantId),
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: fmt.Sprintf("User not found for tenant %s and user ID %s", tenantId, userId.String()),
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("Failed to retrieve user details for tenant %s and user ID %s: %s", tenantId, userId.String(), err.Error()),
			}
		}
	}
	if !userDetails.Status {
		return nil, &dbmodels.ServiceResponse{
			Code:    400,
			Message: fmt.Sprintf("User %s is already disabled for tenant %s", userId.String(), tenantId),
		}
	}
	// Disable the user
	err = u.UserRepo.ChangeStatus(false, userDetails.Id)
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to disable user %s for tenant %s: %s", userId.String(), tenantId, err.Error()),
		}
	}
	return &models.DisableUserResponse{
		Message: fmt.Sprintf("User %s successfully disabled for tenant %s", userId.String(), tenantId),
	}, nil
}

func (u *User) GetUserRole(ctx context.Context, userId uuid.UUID) (*models.GetRoleDetailsUser, error) {
	tenantId := ctx.Value("tenant_id").(string)
	userDetails, err := u.UserRepo.GetUserDetails(dbmodels.DBUser{
		TenantId: uuid.MustParse(tenantId),
		Id:       userId,
	})
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: fmt.Sprintf("User not found for tenant %s and user ID %s", tenantId, userId.String()),
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    500,
				Message: fmt.Sprintf("Failed to retrieve user details for tenant %s and user ID %s: %s", tenantId, userId.String(), err.Error()),
			}
		}
	}
	return &models.GetRoleDetailsUser{
		UserId: userDetails.Id.String(),
		Email:  userDetails.Email,
		Roles:  userDetails.Roles,
	}, nil
}
