package interfaces

import "firebase.google.com/go/v4/messaging"

type Messaging interface {
	GetClient() *messaging.Client
	SendMessage(registrationToken string, data map[string]string) error
}
