package environment

import (
	"be-tpms/src/api/usecases/interfaces"
	"github.com/go-resty/resty/v2"
)

type Env struct {
	RestClient        resty.Client
	CVModelRestClient interfaces.CVModelRestClient
	UserPersister     interfaces.UserPersister
	DogPersister      interfaces.DogPersiter
}
