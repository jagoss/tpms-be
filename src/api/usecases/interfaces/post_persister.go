package interfaces

import "be-tpms/src/api/domain/model"

type PostPersister interface {
	Insert(post *model.Post) (*model.Post, error)
	Get(id int64) (*model.Post, error)
	Update(post *model.Post) (*model.Post, error)
	Delete(id int64) (bool, error)
	DeleteByDogId(id int64) (bool, error)
	GetAll() ([]model.Post, error)
}
