package lostandfound

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
	"math"
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

func (l *LostFoundDogs) GetAllMissingDogsList() []model.Dog {
	return l.dogPersister.GetMissingDogs()
}

func (l *LostFoundDogs) GetMissingDogsInRadius(userLat float64, userLng float64, radius float64) []model.Dog {
	missingDogs := l.dogPersister.GetMissingDogs()
	var dogsInRadio []model.Dog
	for _, dog := range missingDogs {
		if distance(userLat, userLng, float64(dog.Latitude), float64(dog.Longitude)) <= radius {
			dogsInRadio = append(dogsInRadio, dog)
		}
	}
	return dogsInRadio
}

func distance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	radlat1 := math.Pi * lat1 / 180
	radlat2 := math.Pi * lat2 / 180

	theta := lng1 - lng2
	radtheta := math.Pi * theta / 180

	dist := math.Sin(radlat1)*math.Sin(radlat2) + math.Cos(radlat1)*math.Cos(radlat2)*math.Cos(radtheta)
	if dist > 1 {
		dist = 1
	}

	dist = math.Acos(dist)
	dist = dist * 180 / math.Pi
	dist = dist * 60 * 1.1515
	dist = dist * 1.609344
	return dist
}

func (l *LostFoundDogs) PossibleMatchingDogs(dogID uint, matchingDogIDs []uint, userManager interfaces.UserManager, messaging interfaces.Messaging) error {
	for _, id := range matchingDogIDs {
		dog, err := l.dogPersister.GetDog(id)
		if err != nil {
			return err
		}

		data := map[string]string{
			"title": fmt.Sprintf("Puede que alguien viera a %s!", dog.Name),
			"body":  fmt.Sprintf("Confirma la imagen para ver si es %s", dog.Name),
		}
		if err := userManager.SendPushToOwner(dog.Owner.Email, data, messaging); err != nil {
			return err
		}
	}
	return nil
}
