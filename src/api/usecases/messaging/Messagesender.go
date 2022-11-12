package messaging

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
	"log"
)

type MessageSender struct {
	messaging     interfaces.Messaging
	userPersister interfaces.UserPersister
}

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
	user := processUsers(possibleUsers)
	ms.sendMessage(user, reporterUserID, data)
	return nil
}

func processUsers(users []model.User) <-chan model.User {
	out := make(chan model.User)
	go func() {
		for _, user := range users {
			out <- user
		}
		close(out)
	}()
	return out
}

func (ms *MessageSender) sendMessage(user <-chan model.User, reporterID string, data map[string]string) {
	go func() {
		for u := range user {
			if u.ID != reporterID {
				err := ms.messaging.SendMessage(u.FCMToken, data)
				if err != nil {
					log.Printf("error sending message to user %s: %s", u.ID, err.Error())
				}
			}
		}
	}()
}
