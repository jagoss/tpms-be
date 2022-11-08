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

func (dp *DogPersister) InsertDog(dog *model.Dog) (*model.Dog, error) {
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
	tx := dp.db.Connection.Preload("Owner").Preload("Host").First(&dog, dogID)
	if tx.Error != nil {
		if IsRecordNotFoundError(tx.Error) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return &dog, nil
}

func (dp *DogPersister) GetDogs(ids []uint) ([]model.Dog, error) {
	var dogs []model.Dog
	tx := dp.db.Connection.Preload("Owner").Preload("Host").Find(&dogs, ids)
	if tx.Error != nil {
		if IsRecordNotFoundError(tx.Error) {
			return make([]model.Dog, 0), nil
		}
		return nil, tx.Error
	}
	return dogs, nil
}

func (dp *DogPersister) UpdateDog(dog *model.Dog) (*model.Dog, error) {
	tx := dp.db.Connection.Save(dog)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return dog, nil
}

func (dp *DogPersister) DeleteDog(dogID uint) error {
	tx := dp.db.Connection.Preload("Owner").Preload("Host").Delete(&model.Dog{}, dogID)
	return tx.Error
}

func (dp *DogPersister) DogExisitsByNameAndOwner(dogName string, ownerID string) (bool, error) {
	var dog model.Dog
	tx := dp.db.Connection.Preload("Owner").Preload("Host").Where("name = ? AND owner_id = ?", dogName, ownerID).First(&dog)
	if tx.Error != nil {
		if IsRecordNotFoundError(tx.Error) {
			return false, nil
		}
		return false, tx.Error
	}
	return dog.ID != 0, nil
}

func (dp *DogPersister) GetMissingDogs() ([]model.Dog, error) {
	var dogs []model.Dog
	tx := dp.db.Connection.Preload("Owner").Preload("Host").Where("is_lost = ?", "true").Find(&dogs)
	if tx.Error != nil {
		if IsRecordNotFoundError(tx.Error) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return dogs, nil
}

func (dp *DogPersister) GetDogsByUser(userID string) ([]model.Dog, error) {
	var dogs []model.Dog
	tx := dp.db.Connection.Preload("Owner").Preload("Host").Where("host_id = ?", userID).Find(dogs)
	if tx.Error != nil {
		if IsRecordNotFoundError(tx.Error) {
			return nil, nil
		}
		return nil, tx.Error
	}
	return dogs, nil
}
