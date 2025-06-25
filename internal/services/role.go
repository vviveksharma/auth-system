package services

import (
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/models"
	"github.com/vviveksharma/auth/internal/repo"
	dbmodels "github.com/vviveksharma/auth/models"
)

type RoleService interface {
	ListAllRoles() (response []*models.ListAllRolesResponse, err error)
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
