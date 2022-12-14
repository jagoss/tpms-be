package posts

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"log"
)

type PostManager struct {
	postPersister interfaces.PostPersister
}

func NewPostManager(persister interfaces.PostPersister) *PostManager {
	return &PostManager{persister}
}

func (p *PostManager) RegisterPost(post *model.Post) (*model.Post, error) {
	post, err := p.postPersister.Insert(post)
	if err != nil {
		log.Printf("[postmanager.RegisterPost] error persisting post: %s", err.Error())
		return nil, err
	}
	return post, nil
}

func (p *PostManager) RemovePostByDog(dogId int64) (bool, error) {
	return p.postPersister.DeleteByDogId(dogId)
}

func (p *PostManager) GetPost(id int64) (*model.Post, error) {
	return p.postPersister.Get(id)
}

func (p *PostManager) GetAllPost() ([]model.Post, error) {
	return p.postPersister.GetAll()
}
