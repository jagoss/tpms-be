package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"database/sql"
	"fmt"
	"log"
)

type DogPersister struct {
	connection *db.Connection
}

func NewDogPersister(connection *db.Connection) *DogPersister {
	return &DogPersister{connection}
}

func (dp *DogPersister) InsertDog(dog *model.Dog) (*model.Dog, error) {
	query := "INSERT INTO dogs(name, breed, age, size, coat_color, coat_length, is_lost, owner_id, host_id, latitude, longitude, img_url) VALUES (?, ?, ?, ?, ?)"
	result, err := dp.connection.DB.Exec(query, dog.Name, dog.Breed, dog.Age, dog.Size, dog.CoatColor, dog.CoatLength, dog.IsLost, dog.Owner.ID, dog.Host.ID, dog.Latitude, dog.Longitude, dog.ImgUrl)
	if err != nil {
		return nil, err
	}
	if rows, _ := result.RowsAffected(); rows == 0 {
		return nil, fmt.Errorf("dog %v was not inserted", dog)
	}

	dog.ID, _ = result.LastInsertId()
	return dog, nil
}

func (dp *DogPersister) GetDog(dogID uint) (*model.Dog, error) {
	query := "SELECT * FROM tpms_prod.dogs WHERE id = ?"
	row := dp.connection.DB.QueryRow(query, dogID)
	if row.Err() != nil {
		return nil, row.Err()
	}

	var dog model.DogModel
	if err := row.Scan(&dog.ID, &dog.Name, &dog.Breed, &dog.Age, &dog.Size, &dog.CoatColor, &dog.CoatLength, &dog.IsLost, &dog.OwnerID, &dog.HostID, &dog.Latitude, &dog.Longitude, &dog.ImgUrl); err != nil {
		return nil, err
	}
	up := UserPersister{dp.connection}
	var owner, host *model.User
	if dog.OwnerID != "" {
		owner, _ = up.GetUser(dog.OwnerID)
	}
	if dog.HostID != "" {
		host, _ = up.GetUser(dog.HostID)
	}

	resultDog := mapToDog(dog, owner, host)

	return &resultDog, nil
}

func (dp *DogPersister) GetDogs(ids []uint) ([]model.Dog, error) {
	query := "SELECT * FROM tpms_prod.dogs WHERE id in (?)"
	rows, err := dp.connection.DB.Query(query, ids)
	if err != nil {
		return nil, err
	}
	resultList, err := dp.parseDogs(rows)

	if err != nil {
		return nil, err
	}

	return resultList, nil
}

func (dp *DogPersister) UpdateDog(dog *model.Dog) (*model.Dog, error) {
	query := "UPDATE tpms_prod.dogs SET name = ?, age = ?, breed = ?, size = ?, coat_color=?, coat_length = ?, is_lost = ?, latitude = ?, longitude = ?, img_url = ? WHERE id = ?"
	result, err := dp.connection.DB.Exec(query, dog.Name, dog.Age, dog.Breed, dog.Size, dog.CoatColor, dog.CoatLength, dog.IsLost, dog.Latitude, dog.Longitude, dog.ImgUrl)
	if err != nil {
		return nil, err
	}
	if amount, err := result.RowsAffected(); err != nil || amount == 0 {
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("no rows affected in database when updating dog %v", dog)
	}

	return dog, nil
}

func (dp *DogPersister) DeleteDog(dogID uint) error {
	query := "DELETE FROM dogs WHERE id = ?"
	exec, err := dp.connection.DB.Exec(query, dogID)
	if err != nil {
		return err
	}
	if count, _ := exec.RowsAffected(); count == 0 {
		return fmt.Errorf("dog %d was not deleted", dogID)
	}
	return nil
}

func (dp *DogPersister) DogExisitsByNameAndOwner(dogName string, ownerID string) (bool, error) {
	query := "SELECT * FROM dogs WHERE name = ? AND owner_id = ?"
	rows, err := dp.connection.DB.Query(query, dogName, ownerID)
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}

func (dp *DogPersister) GetMissingDogs() ([]model.Dog, error) {
	query := "SELECT * FROM dogs where is_lost = true AND owner_id IS NULL"
	rows, err := dp.connection.DB.Query(query)
	if err != nil {
		return nil, err
	}

	dogs, err := dp.parseDogs(rows)
	if err != nil {
		return nil, err
	}

	return dogs, nil
}

func (dp *DogPersister) GetDogsByUser(userID string) ([]model.Dog, error) {
	query := "SELECT * FROM dogs WHERE host_id = ? OR owner_id = ?"

	rows, err := dp.connection.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}

	dogs, err := dp.parseDogs(rows)
	if err != nil {
		return nil, err
	}
	log.Printf("dogs %v", dogs)
	return dogs, nil
}

func mapToDog(dogModel model.DogModel, owner *model.User, host *model.User) model.Dog {
	return model.Dog{ID: dogModel.ID,
		Name:       dogModel.Name,
		Breed:      dogModel.Breed,
		Size:       dogModel.Size,
		Age:        dogModel.Age,
		CoatColor:  dogModel.CoatColor,
		CoatLength: dogModel.CoatLength,
		IsLost:     dogModel.IsLost,
		Owner:      owner,
		Host:       host,
		Latitude:   dogModel.Latitude,
		Longitude:  dogModel.Longitude,
		ImgUrl:     dogModel.ImgUrl,
	}
}

func (dp *DogPersister) parseDogs(rows *sql.Rows) ([]model.Dog, error) {
	var resultList []model.Dog
	up := UserPersister{dp.connection}

	for rows.Next() {
		var dog model.DogModel
		if err := rows.Scan(&dog.ID, &dog.Name, &dog.Breed, &dog.Age, &dog.Size, &dog.CoatColor, &dog.CoatLength, &dog.IsLost, &dog.OwnerID, &dog.HostID, &dog.Latitude, &dog.Longitude, &dog.ImgUrl); err != nil {
			return nil, err
		}
		var owner, host *model.User
		if dog.OwnerID != "" {
			owner, _ = up.GetUser(dog.OwnerID)
		}
		if dog.HostID != "" {
			host, _ = up.GetUser(dog.HostID)
		}

		resultList = append(resultList, mapToDog(dog, owner, host))
	}
	if resultList == nil || len(resultList) == 0 {
		return make([]model.Dog, 0), nil
	}
	return resultList, nil
}
