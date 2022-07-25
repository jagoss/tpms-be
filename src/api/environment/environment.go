package environment

import (
	"be-tpms/src/api/usecases/interfaces"
	"firebase.google.com/go/auth"
	"github.com/go-resty/resty/v2"
)

type Env struct {
	FirebaseAuth      auth.Client
	RestClient        resty.Client
	CVModelRestClient interfaces.CVModelRestClient
	UserPersister     interfaces.UserPersister
	DogPersister      interfaces.DogPersister
}

func InitEnv(firebaseAuth auth.Client, restClient resty.Client, cvModelClient interfaces.CVModelRestClient, userPersister interfaces.UserPersister, dogPersister interfaces.DogPersister) *Env {
	return &Env{
		FirebaseAuth:      firebaseAuth,
		RestClient:        restClient,
		CVModelRestClient: cvModelClient,
		UserPersister:     userPersister,
		DogPersister:      dogPersister,
	}
}
