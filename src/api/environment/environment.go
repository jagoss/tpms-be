package environment

import (
	"be-tpms/src/api/usecases/interfaces"
	"github.com/go-resty/resty/v2"
)

type Env struct {
	RestClient        resty.Client
	CVModelRestClient interfaces.CVModelRestClient
	UserPersister     interfaces.UserPersister
	DogPersister      interfaces.DogPersister
}

func InitEnv(restClient resty.Client, cvModelClient interfaces.CVModelRestClient, userPersister interfaces.UserPersister, dogPersister interfaces.DogPersister) *Env {
	return &Env{
		RestClient:        restClient,
		CVModelRestClient: cvModelClient,
		UserPersister:     userPersister,
		DogPersister:      dogPersister,
	}
}
