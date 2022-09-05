package interfaces

import "be-tpms/src/api/domain/model"

type DogManager interface {
	Get(uint) (*model.Dog, error)
	Register(*model.Dog, [][]byte, UserManager) (*model.Dog, error)
	Modify(*model.Dog, [][]byte, UserManager) (*model.Dog, error)
	Delete(uint) (bool, error)
}
