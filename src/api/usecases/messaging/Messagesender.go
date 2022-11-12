package messaging

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
	"log"
)

type MessageSender struct {
	messaging     interfaces.Messaging
	userPersister interfaces.UserPersister
}

const radius = 2

func NewMessageSender(messaging interfaces.Messaging, persister interfaces.UserPersister) *MessageSender {
	return &MessageSender{
		messaging:     messaging,
		userPersister: persister,
	}
}

func (ms *MessageSender) SendToEnabledUsers(dog *model.Dog) error {
	data := map[string]string{
		io.TITLE: "Se ha perdido un perro cerca tuyo!",
		io.BODY:  fmt.Sprintf("Se perdio un perro de raza %s cerca tuyo!", dog.Breed.String()),
	}

	var reporterUserID string
	if dog.Owner == nil {
		reporterUserID = dog.Host.ID
	} else {
		reporterUserID = dog.Owner.ID
	}

	possibleUsers, err := ms.userPersister.GetUsersEnabledMessages()
	if err != nil {
		return err
	}

	user := processUsers(dog, possibleUsers, reporterUserID)
	ms.sendMessage(user, data)
	return nil
}

func processUsers(dog *model.Dog, users []model.User, reporterID string) <-chan model.User {
	out := make(chan model.User)
	dogLat := dog.Latitude
	dogLong := dog.Longitude
	go func() {
		for _, user := range users {
			if user.ID != reporterID && usecases.Distance(dogLat, dogLong, user.Latitude, user.Longitude) < radius {
				out <- user
			}
		}
		close(out)
	}()
	return out
}

func (ms *MessageSender) sendMessage(user <-chan model.User, data map[string]string) {
	go func() {
		for u := range user {
			err := ms.messaging.SendMessage(u.FCMToken, data)
			if err != nil {
				log.Printf("error sending message to user %s: %s", u.ID, err.Error())
			}
		}
	}()
}
