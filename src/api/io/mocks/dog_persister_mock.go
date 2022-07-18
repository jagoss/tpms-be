package mock

import (
	"be-tpms/src/api/domain/model"
)

type DogPersisterMock struct {
	dogs          []model.Dog
	getCounter    int
	errors        []error
	errorsCounter int
}

func (dp *DogPersisterMock) InsertDog(dog *model.Dog) (*model.Dog, error) {
	d := dp.dogs[dp.getCounter]
	dp.getCounter += 1
	return &d, nil
}

func (dp *DogPersisterMock) GetDog(dogID uint) (*model.Dog, error) {
	dog := dp.dogs[dp.getCounter]
	dp.getCounter += 1
	return &dog, nil
}

func (dp *DogPersisterMock) UpdateDog(dog *model.Dog) (*model.Dog, error) {
	return dog, nil
}

func (dp *DogPersisterMock) DeleteDog(dogID uint) error {
	err := dp.errors[dp.errorsCounter]
	dp.errorsCounter += 1
	return err
}
