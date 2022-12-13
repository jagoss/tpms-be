package interfaces

import "be-tpms/src/api/domain/model"

type PostManager interface {
	RegisterPost(model *model.Post) (*model.Post, error)
	RemovePostByDog(dogId int64) (bool, error)
	GetPost(id int64) (*model.Post, error)
	GetAllPost() ([]model.Post, error)
}
