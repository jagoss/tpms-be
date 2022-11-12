package interfaces

import "be-tpms/src/api/domain/model"

type UserPersister interface {
	GetUser(string) (*model.User, error)
	InsertUser(*model.User) (*model.User, error)
	UpdateUser(*model.User) (*model.User, error)
	DeleteUser(string) error
	GetUsersEnabledMessages() ([]model.User, error)
}
