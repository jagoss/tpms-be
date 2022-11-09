package lostandfound

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
	"log"
	"math"
)

type LostFoundDogs struct {
	dogPersister           interfaces.DogPersister
	userPersister          interfaces.UserPersister
	possibleMatchPersister interfaces.PossibleMatchPersister
}

func NewLostFoundDogs(dogPersister interfaces.DogPersister, userPersister interfaces.UserPersister, possibleMatchPersister interfaces.PossibleMatchPersister) LostFoundDogs {
	return LostFoundDogs{
		dogPersister:           dogPersister,
		userPersister:          userPersister,
		possibleMatchPersister: possibleMatchPersister,
	}
}

func (l *LostFoundDogs) ReuniteDog(dogID uint, matchingDogID uint, sender interfaces.Messaging) (*model.Dog, error) {
	dog, err := l.dogPersister.GetDog(dogID)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] %v", err)
	}
	possibleDog, err := l.dogPersister.GetDog(matchingDogID)
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.ReuniteDog] %v", err)
	}
	err = l.unifyDogs(dog, possibleDog)
	if err != nil {
		return nil, err
	}

	dogIDsRemoved, err := l.cleanPossibleMatches(dogID)
	if err != nil {
		return nil, err
	}

	err = l.notifyDogsHosters(dog.Name, dogIDsRemoved, sender)
	if err != nil {
		return nil, err
	}

	return dog, nil
}

func (l *LostFoundDogs) GetAllMissingDogsList() ([]model.Dog, error) {
	return l.dogPersister.GetMissingDogs()
}

func (l *LostFoundDogs) GetMissingDogsInRadius(userLat float64, userLng float64, radius float64) ([]model.Dog, error) {
	missingDogs, err := l.dogPersister.GetMissingDogs()
	if err != nil {
		return nil, fmt.Errorf("[lostfounddogs.GetMissingDogsInRadius] error getting dogs: %s", err.Error())
	}
	var dogsInRadio []model.Dog
	for _, dog := range missingDogs {
		if distance(userLat, userLng, float64(dog.Latitude), float64(dog.Longitude)) <= radius {
			dogsInRadio = append(dogsInRadio, dog)
		}
	}
	return dogsInRadio, nil
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

func (l *LostFoundDogs) PossibleMatchingDogs(dogID uint, possibleDogsIDs []uint, sender interfaces.Messaging) error {
	for _, id := range possibleDogsIDs {
		dog, err := l.dogPersister.GetDog(id)
		if err != nil {
			return err
		}

		if err = l.possibleMatchPersister.AddPossibleMatch(dogID, id); err != nil {
			log.Printf("%v", err)
		} else {
			data := map[string]string{
				"title": fmt.Sprintf("Puede que alguien viera a %s!", dog.Name),
				"body":  fmt.Sprintf("Confirma la imagen para ver si es %s", dog.Name),
			}

			if err = sender.SendMessage(dog.Owner.FCMToken, data); err != nil {
				log.Printf("error sending push notification to user %s: %v", dog.Owner.ID, err)
			}
		}
	}
	return nil
}

func (l *LostFoundDogs) AcknowledgePossibleDog(dogID uint, possibleDogID uint, sender interfaces.Messaging) error {
	possibleDog, _ := l.dogPersister.GetDog(possibleDogID)
	dog, _ := l.dogPersister.GetDog(dogID)

	data := map[string]string{
		"title": fmt.Sprintf("Han confirmado tu perro!"),
		"body":  fmt.Sprintf("Puede que %s tenga a %s", possibleDog.Owner.Name, dog.Name),
	}
	return l.updatePossibleDogMatch(dogID, possibleDogID, model.Accepted, data, sender)
}

func (l *LostFoundDogs) RejectPossibleDog(dogID uint, possibleDogID uint, sender interfaces.Messaging) error {
	possibleDog, _ := l.dogPersister.GetDog(possibleDogID)
	dog, _ := l.dogPersister.GetDog(dogID)

	data := map[string]string{
		"title": fmt.Sprintf("Han rechazado tu match"),
		"body":  fmt.Sprintf("Parece que %s no era %s. Sigamos buscando!", possibleDog.Name, dog.Name),
	}
	return l.updatePossibleDogMatch(dogID, possibleDogID, model.Accepted, data, sender)
}

func (l *LostFoundDogs) updatePossibleDogMatch(dogID uint, possibleDogID uint, ack model.Ack, data map[string]string, sender interfaces.Messaging) error {
	if ack == model.Accepted {
		err := l.possibleMatchPersister.UpdateAck(dogID, possibleDogID, ack)
		if err != nil {
			return fmt.Errorf("[LostFoundDogs.updatePossibleDogMatch] error updating ack dog %d wiht dog %d: %v",
				dogID, possibleDogID, err)
		}
	} else {
		err := l.possibleMatchPersister.Delete(dogID, possibleDogID)
		if err != nil {
			return fmt.Errorf("[LostFoundDogs.updatePossibleDogMatch] error deleting possible match with dogs dog %d wiht dog %d: %v",
				dogID, possibleDogID, err)
		}
	}

	dog, _ := l.dogPersister.GetDog(dogID)

	if err := sender.SendMessage(dog.Owner.FCMToken, data); err != nil {
		log.Printf("error sending push notification to user %s: %v", dog.Owner.ID, err)
	}

	return nil
}

func (l *LostFoundDogs) unifyDogs(dog *model.Dog, matchingDog *model.Dog) error {
	dog.IsLost = false
	dog.Host = dog.Owner
	dog.ImgUrl = fmt.Sprintf("%s;%s", dog.ImgUrl, matchingDog.ImgUrl)
	notifyModel(dog.ID, matchingDog.ID)
	return nil
}

func (l *LostFoundDogs) cleanPossibleMatches(dogID uint) ([]uint, error) {
	removedPossibleDogs, err := l.possibleMatchPersister.RemovePossibleDogMatches(dogID)
	if err != nil {
		return nil, err
	}

	removedDogPossibleMatches, err := l.possibleMatchPersister.RemovePossibleMatchesForDog(dogID)
	if err != nil {
		return nil, err
	}

	possibleMatches := append(removedPossibleDogs, removedDogPossibleMatches...)
	var ids []uint
	for _, dog := range possibleMatches {
		if dog.DogID != 0 {
			ids = append(ids, dog.DogID)
		} else {
			ids = append(ids, dog.PossibleDogID)
		}
	}

	return ids, nil
}

func (l *LostFoundDogs) notifyDogsHosters(actualDogName string, removedDogs []uint, sender interfaces.Messaging) error {
	dogs, err := l.dogPersister.GetDogs(removedDogs)
	if err != nil {
		return err
	}

	data := map[string]string{
		"title": fmt.Sprintf("Parece que era otro perro!"),
	}
	for _, dog := range dogs {
		err = sender.SendMessage(dog.Host.FCMToken, data)
		data["body"] = fmt.Sprintf("Parece que %s no era %s. Sigamos buscando!", dog.Name, actualDogName)
		if err != nil {
			log.Printf("error sending notification tu user %s: %s", dog.Host.ID, err.Error())
		}
	}

	return nil
}

func notifyModel(dogID int64, sameDogID int64) {

}
