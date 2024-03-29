package environment

import (
	"be-tpms/src/api/usecases/interfaces"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"net/http"
)

type Env struct {
	FirebaseApp            *firebase.App
	FirebaseAuth           auth.Client
	NotificationSender     interfaces.Messaging
	RestClient             *http.Client
	CVModelRestClient      interfaces.CVModelRestClient
	UserPersister          interfaces.UserPersister
	DogPersister           interfaces.DogPersister
	PossibleMatchPersister interfaces.PossibleMatchPersister
	PostPersister          interfaces.PostPersister
	Storage                interfaces.Storage
}

func InitEnv(firebaseApp *firebase.App, firebaseAuth auth.Client, notificationSender interfaces.Messaging,
	restClient *http.Client, cvModelClient interfaces.CVModelRestClient, userPersister interfaces.UserPersister,
	dogPersister interfaces.DogPersister, possibleMatchPersister interfaces.PossibleMatchPersister,
	postPersister interfaces.PostPersister, storage interfaces.Storage) *Env {
	return &Env{
		FirebaseApp:            firebaseApp,
		FirebaseAuth:           firebaseAuth,
		NotificationSender:     notificationSender,
		RestClient:             restClient,
		CVModelRestClient:      cvModelClient,
		UserPersister:          userPersister,
		DogPersister:           dogPersister,
		PossibleMatchPersister: possibleMatchPersister,
		PostPersister:          postPersister,
		Storage:                storage,
	}
}
