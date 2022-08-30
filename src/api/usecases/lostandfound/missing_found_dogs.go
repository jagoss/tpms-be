package lostandfound

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
)

type LostFoundDogs struct {
	dogPersister  interfaces.DogPersister
	userPersister interfaces.UserPersister
}

func NewLostFoundDogs(dogPersister interfaces.DogPersister, userPersister interfaces.UserPersister) LostFoundDogs {
	return LostFoundDogs{
		dogPersister:  dogPersister,
		userPersister: userPersister,
	}
}

func (l *LostFoundDogs) ReuniteDog(dogID uint, ownerID string, hosterID string) (*model.Dog, error) {
	if ownerID == hosterID {
		dog, _ := l.dogPersister.GetDog(dogID)
		return dog, nil
	}
	dog, err := l.dogPersister.GetDog(dogID)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] %v", err)
	}
	owner, err := l.userPersister.GetUser(ownerID)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] %v", err)
	}
	dog.Host = owner
	dog.IsLost = false
	modifiedDog, err := l.dogPersister.UpdateDog(dog)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] error updating owner: %v", err)
	}
	if modifiedDog.IsLost || modifiedDog.Owner.ID != ownerID {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] error updating dog:\n correct ownerID: %s, but is %s and it may remains mark as lost!", ownerID, modifiedDog.Owner.ID)
	}

	return modifiedDog, nil
}

func (l *LostFoundDogs) GetMissingDogsList() []model.Dog {
	return l.dogPersister.GetMissingDogs()
}
