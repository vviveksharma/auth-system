package services

import (
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type RoleService interface {
	CreateCustomRole(roleName string) (resp *models.CreateCustomRoleResponse, err error)
	ListAllRoles() (response []*models.ListAllRolesResponse, err error)
	VerifyRole(req *models.VerifyRoleRequest) (response *models.VerifyRoleResponse, err error)
}

type Role struct {
	RoleRepo repo.RoleRepositoryInterface
}

func NewRoleService() (RoleService, error) {
	ser := &Role{}
	err := ser.SetupRepo()
	if err != nil {
		return nil, err
	}
	return ser, nil
}

func (r *Role) SetupRepo() error {
	var err error
	role, err := repo.NewRoleRepository(db.DB)
	if err != nil {
		return err
	}
	r.RoleRepo = role
	return nil
}

func (r *Role) ListAllRoles() (response []*models.ListAllRolesResponse, err error) {
	roles, err := r.RoleRepo.GetAllRoles()
	if err != nil {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while listing all the roles: " + err.Error(),
		}
	}
	for _, i := range roles {
		var res models.ListAllRolesResponse
		res.Name = i.Role
		response = append(response, &res)
	}
	return response, nil
}

func (r *Role) VerifyRole(req *models.VerifyRoleRequest) (response *models.VerifyRoleResponse, err error) {
	roleDetails, err := r.RoleRepo.FindRoleId(req.RoleName)
	if err != nil {
		if err.Error() == "record not found" {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "No role with this name exists",
			}
		} else {
			return nil, &dbmodels.ServiceResponse{
				Code:    404,
				Message: "error while finding the role:" + err.Error(),
			}
		}
	}
	log.Info("the roleId: ", roleDetails)
	if roleDetails != uuid.MustParse(req.RoleId) {
		return nil, &dbmodels.ServiceResponse{
			Code:    404,
			Message: "Role name and ID do not match. Please provide the correct role ID and role name.",
		}
	}
	return &models.VerifyRoleResponse{
		Message: true,
	}, nil
}

func (r *Role) CreateCustomRole(roleName string) (resp *models.CreateCustomRoleResponse, err error) {
	_, err = r.RoleRepo.FindRoleId(roleName)
	if err != nil && err.Error() != "record not found " {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while fetching the existing role information: " + err.Error(),
		}
	}
	err = r.RoleRepo.CreateRole(&dbmodels.DBRoles{
		Role:     roleName,
		RoleId:   uuid.New(),
		RoleType: "custom",
	})
	if err != nil && err.Error() != "record not found" {
		return nil, &dbmodels.ServiceResponse{
			Code:    500,
			Message: "error while creating a role: " + err.Error(),
		}
	}
	return &models.CreateCustomRoleResponse{
		Message: "role with " + roleName + " created successfully.",
	}, nil
}
