package lostandfound

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/dogs"
	"be-tpms/src/api/usecases/interfaces"
	"be-tpms/src/api/usecases/users"
	"fmt"
)

type LostFoundDogs struct {
	dogManager  interfaces.DogManager
	userManager interfaces.UserManager
}

func NewLostFoundDogs(dogPersister interfaces.DogPersister, userPersister interfaces.UserPersister) LostFoundDogs {
	return LostFoundDogs{
		dogManager:  dogs.NewDogManager(dogPersister),
		userManager: users.NewUserManager(userPersister),
	}
}

func (l *LostFoundDogs) ReuniteDog(dogID uint, ownerID string, hosterID string) (*model.Dog, error) {
	if ownerID == hosterID {
		dog, _ := l.dogManager.Get(dogID)
		return dog, nil
	}
	dog, err := l.dogManager.Get(dogID)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] %v", err)
	}
	owner, err := l.userManager.Get(ownerID)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] %v", err)
	}
	dog.Host = owner
	dog.IsLost = false
	modifiedDog, err := l.dogManager.Modify(dog)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] error updating owner: %v", err)
	}
	if modifiedDog.IsLost || modifiedDog.Owner.ID != ownerID {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] error updating dog:\n correct ownerID: %s, but is %s and it may remains mark as lost!", ownerID, modifiedDog.Owner.ID)
	}

	return modifiedDog, nil
}
