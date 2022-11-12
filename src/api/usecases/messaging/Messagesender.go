package messaging

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
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

func (ms *MessageSender) SendToEnabledUsers(dog *model.Dog, data map[string]string) error {
	return nil
}
