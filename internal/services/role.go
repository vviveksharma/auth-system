package services

import (
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/repo"
)

type RoleService interface{}

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
