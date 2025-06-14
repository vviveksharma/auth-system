package services

import (
	"github.com/vviveksharma/auth/db"
	"github.com/vviveksharma/auth/internal/repo"
)

type IAuth interface{}

type Auth struct {
	UserRepo repo.UserRepositoryInterface
}

func NewAuthService() (IAuth, error) {
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
	return nil
}

func (a *Auth) Login()