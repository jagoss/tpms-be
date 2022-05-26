package usermanager

import "github.com/tesis/be-tpms/src/api/usecases/interfaces"

type UserManager struct {
	userRepo *interfaces.UserRepository
}

func NewUserManager(userRepo *interfaces.UserRepository) *UserManager {
	return &UserManager{userRepo: userRepo}
}
