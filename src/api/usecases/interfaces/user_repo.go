package interfaces

import "github.com/tesis/be-tpms/src/api/domain/model"

type UserRepository interface {
	Read(string) (*model.User, error)
	Create(*model.User) (*model.User, error)
	Update(*model.User) (*model.User, error)
	Delete(string) (bool, error)
}
