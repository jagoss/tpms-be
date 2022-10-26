package interfaces

import "be-tpms/src/api/domain/model"

type UserManager interface {
	Get(string) (*model.User, error)
	Register(*model.User) (*model.User, error)
	Modify(*model.User) (*model.User, error)
	Delete(string) (bool, error)
	UpdateFCMToken(id string, token string) error
	SendPushToOwner(string, map[string]string, Messaging) error
}
