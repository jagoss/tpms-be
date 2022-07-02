package users

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
)

type UserManager struct {
	userRepo *interfaces.UserPersister
}

func NewUserManager(userRepo *interfaces.UserPersister) *UserManager {
	return &UserManager{userRepo: userRepo}
}

func (u *UserManager) Get(string) (*model.User, error) {
	return nil, nil
}

func (u *UserManager) Register(*model.User) (*model.User, error) {

	return nil, nil
}

func (u *UserManager) Modify(*model.User) (*model.User, error) {
	return nil, nil
}
func (u *UserManager) Delete(string) (bool, error) {
	return true, nil
}
