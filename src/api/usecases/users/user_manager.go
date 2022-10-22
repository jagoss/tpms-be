package users

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
)

type UserManager struct {
	userRepo interfaces.UserPersister
}

func NewUserManager(userRepo interfaces.UserPersister) *UserManager {
	return &UserManager{userRepo: userRepo}
}

func (u *UserManager) Get(userID string) (*model.User, error) {
	user, err := u.userRepo.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("[usermanager.Get] error getting user with id %s: %v", userID, err)
	}
	return user, nil
}

func (u *UserManager) Register(user *model.User) (*model.User, error) {
	user, err := u.userRepo.InsertUser(user)
	if err != nil {
		return nil, fmt.Errorf("[usermanager.Register] error registing user: %w", err)
	}
	return user, nil
}

func (u *UserManager) Modify(user *model.User) (*model.User, error) {
	user, err := u.userRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("[usermanager.Modify] error modifying user with id %s: %v", user.ID, err)
	}
	return user, nil
}

func (u *UserManager) Delete(userID string) (bool, error) {
	if err := u.userRepo.DeleteUser(userID); err != nil {
		return false, fmt.Errorf("[usermanager.Delete] error registing user with id %s: %v", userID, err)
	}
	return true, nil
}

func (u *UserManager) UpdateFCMToken(id string, token string) error {
	user, err := u.userRepo.GetUser(id)
	if err != nil {
		return fmt.Errorf("error getting user %s", id)
	}
	user.FCMToken = token
	user, err = u.userRepo.UpdateUser(user)
	if err != nil {
		return fmt.Errorf("error updating user token for user %s", id)
	}
	return nil
}
