package users

import "be-tpms/src/api/usecases/interfaces"

type UserManager struct {
	userRepo *interfaces.UserPersister
}

func NewUserManager(userRepo *interfaces.UserPersister) *UserManager {
	return &UserManager{userRepo: userRepo}
}
