package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"fmt"
)

type DogPersister struct {
	db *db.DataBase
}

func NewDogPersister(db *db.DataBase) *DogPersister {
	return &DogPersister{db: db}
}

func (dp *DogPersister) InsertDog(dog *model.User) (*model.User, error) {
	result := dp.db.Connection.Create(dog)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no dog was inserted")
	}
	return dog, nil
}

func (dp *DogPersister) GetDog(dogID uint) (*model.Dog, error) {
	var dog model.Dog
	tx := dp.db.Connection.First(&dog, dogID)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return &dog, nil
}

func (dp *DogPersister) UpdateDog(dog *model.Dog) (*model.Dog, error) {
	tx := dp.db.Connection.Save(dog)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dog, nil
}
func (dp *DogPersister) DeleteDog(userID string) error {
	tx := dp.db.Connection.Delete(&model.Dog{}, userID)
	return tx.Error
}
