package lostandfound

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
	"log"
)

type LostFoundDogs struct {
	dogPersister           interfaces.DogPersister
	userPersister          interfaces.UserPersister
	possibleMatchPersister interfaces.PossibleMatchPersister
	postPersister          interfaces.PostPersister
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

	if possibleDog.Owner == nil && possibleDog.Host == nil {
		_, err := l.postPersister.DeleteByDogId(possibleDog.ID)
		if err != nil {
			return nil, err
		}
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
	log.Printf("all missing dogs: %v", missingDogs)
	var dogsInRadio []model.Dog
	for _, dog := range missingDogs {
		if usecases.Distance(userLat, userLng, dog.Latitude, dog.Longitude) <= radius {
			dogsInRadio = append(dogsInRadio, dog)
		}
	}
	log.Printf("Dogs in radius %d: %v", radius, dogsInRadio)
	return dogsInRadio, nil
}

func (l *LostFoundDogs) PossibleMatchingDogs(dogID uint, possibleDogsIDs []uint, sender interfaces.Messaging) error {
	_, err := l.dogPersister.GetDog(dogID)
	if err != nil {
		return fmt.Errorf("[PossibleMatchingDogs] error getting dog with id %d: %s", dogID, err.Error())
	}
	for _, id := range possibleDogsIDs {
		dog, err := l.dogPersister.GetDog(id)
		if err != nil {
			return fmt.Errorf("[PossibleMatchingDogs] error getting possible dog with id %d: %s", id, err.Error())
		}

		if err := l.possibleMatchPersister.AddPossibleMatch(dogID, id); err != nil {
			log.Printf("[PossibleMatchingDogs] %v", err)
		} else {
			send(dog, sender)
		}
	}
	return nil
}

func (l *LostFoundDogs) AcknowledgePossibleDog(dogID uint, possibleDogID uint, sender interfaces.Messaging) error {
	possibleDog, err := l.dogPersister.GetDog(possibleDogID)
	if err != nil {
		return err
	}
	dog, err := l.dogPersister.GetDog(dogID)
	if err != nil {
		return err
	}

	log.Printf("dog: %s; possible dog: %s", dog.Name, possibleDog.Name)

	var user string
	if possibleDog.Owner != nil {
		user = possibleDog.Owner.Name
	} else {
		user = possibleDog.Host.Name
	}

	data := map[string]string{
		"title": fmt.Sprintf("Han confirmado tu perro!"),
		"body":  fmt.Sprintf("Puede que %s tenga a %s", user, dog.Name),
	}
	return l.updatePossibleDogMatch(dogID, possibleDogID, model.Accepted, data, sender)
}

func (l *LostFoundDogs) RejectPossibleDog(dogID uint, possibleDogID uint, sender interfaces.Messaging) error {
	possibleDog, err := l.dogPersister.GetDog(possibleDogID)
	if err != nil {
		return err
	}
	dog, err := l.dogPersister.GetDog(dogID)
	if err != nil {
		return err
	}

	data := map[string]string{
		"title": fmt.Sprintf("Han rechazado tu match"),
		"body":  fmt.Sprintf("Parece que %s no era %s. Sigamos buscando!", possibleDog.Name, dog.Name),
	}
	return l.updatePossibleDogMatch(dogID, possibleDogID, model.Accepted, data, sender)
}

func (l *LostFoundDogs) updatePossibleDogMatch(dogID uint, possibleDogID uint, ack model.Ack, data map[string]string, sender interfaces.Messaging) error {
	if ack == model.Accepted {
		if err := l.possibleMatchPersister.UpdateAck(dogID, possibleDogID, ack); err != nil {
			return fmt.Errorf("[LostFoundDogs.updatePossibleDogMatch] error updating ack dog %d wiht dog %d: %v",
				dogID, possibleDogID, err)
		}
	} else {
		if err := l.possibleMatchPersister.Delete(dogID, possibleDogID); err != nil {
			return fmt.Errorf("[LostFoundDogs.updatePossibleDogMatch] error deleting possible match with dogs dog %d wiht dog %d: %v",
				dogID, possibleDogID, err)
		}
	}

	dog, err := l.dogPersister.GetDog(possibleDogID)
	if err != nil {
		return err
	}

	var userToken string
	if dog.Owner != nil {
		userToken = dog.Owner.FCMToken
	} else {
		userToken = dog.Host.FCMToken
	}

	if err := sender.SendMessage(userToken, data); err != nil {
		log.Printf("error sending push notification to user %s: %v", dog.Owner.ID, err)
	}

	return nil
}

func (l *LostFoundDogs) unifyDogs(dog *model.Dog, matchingDog *model.Dog) error {
	dog.IsLost = false
	dog.Host = dog.Owner
	dog.ImgUrl = fmt.Sprintf("%s;%s", dog.ImgUrl, matchingDog.ImgUrl)
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

func (l *LostFoundDogs) GetPossibleMatchingDogs(id uint, acks []model.Ack) ([]model.PossibleMatch, error) {
	possibleMatches, err := l.possibleMatchPersister.GetPossibleMatches(id, acks)
	if err != nil {
		log.Printf("[LostFoundDogs.GetPossibleMatchingDogs] error getting matches: %s", err.Error())
		return nil, err
	}

	return possibleMatches, nil
}

func send(dog *model.Dog, sender interfaces.Messaging) {
	if dog.Owner != nil {
		data := map[string]string{
			"title": fmt.Sprintf("Puede que alguien viera a %s!", dog.Name),
			"body":  fmt.Sprintf("Confirma la imagen para ver si es %s", dog.Name),
		}
		sendNotification(data, dog.Owner, sender)
	}
	if dog.Host != nil {
		data := map[string]string{
			"title": fmt.Sprintf("Parece que hay un posible dueño de %s!", dog.Name),
			"body":  fmt.Sprintf("Confirma la imagen para ver si es %s", dog.Name),
		}
		sendNotification(data, dog.Host, sender)
	}
}

func sendNotification(data map[string]string, user *model.User, sender interfaces.Messaging) {
	if err := sender.SendMessage(user.FCMToken, data); err != nil {
		log.Printf("error sending push notification to user %s: %v", user.ID, err)
	}
}
