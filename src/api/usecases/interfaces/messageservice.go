package interfaces

import "be-tpms/src/api/domain/model"

type MessageSender interface {
	SendToEnabledUsers(dog *model.Dog) error
}
