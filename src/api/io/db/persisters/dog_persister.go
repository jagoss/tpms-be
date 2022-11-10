package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"database/sql"
	"log"
)

type DogPersister struct {
	connection *db.Connection
}

func NewDogPersister(connection *db.Connection) *DogPersister {
	return &DogPersister{connection}
}

func (dp *DogPersister) InsertDog(dog *model.Dog) (*model.Dog, error) {
	dogModel := mapToDogModel(*dog)

	query := "INSERT INTO dogs(name, breed, age, size, coat_color, coat_length, is_lost, owner_id, host_id, latitude, longitude, img_url) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := dp.connection.DB.Exec(query, dogModel.Name, dogModel.Breed, dogModel.Age, dogModel.Size, dogModel.CoatColor, dogModel.CoatLength, dogModel.IsLost, dogModel.OwnerID, dogModel.HostID, dogModel.Latitude, dogModel.Longitude, dogModel.ImgUrl)
	if err != nil {
		return nil, err
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
	dogModel := mapToDogModel(*dog)
	query := "UPDATE tpms_prod.dogs SET name = ?, age = ?, breed = ?, size = ?, coat_color=?, coat_length = ?, is_lost = ?, latitude = ?, longitude = ?, img_url = ? WHERE id = ?"
	_, err := dp.connection.DB.Exec(query, dogModel.Name, dogModel.Age, dogModel.Breed, dogModel.Size, dogModel.CoatColor, dogModel.CoatLength, dogModel.IsLost, dogModel.Latitude, dogModel.Longitude, dogModel.ImgUrl)
	if err != nil {
		return nil, err
	}

	return dog, nil
}

func (dp *DogPersister) DeleteDog(dogID uint) error {
	query := "DELETE FROM dogs WHERE id = ?"
	_, err := dp.connection.DB.Exec(query, dogID)
	if err != nil {
		return err
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

func mapToDogModel(dog model.Dog) model.DogModel {
	dogModel := model.DogModel{
		Name:       dog.Name,
		Breed:      dog.Breed,
		Size:       dog.Size,
		Age:        dog.Age,
		CoatColor:  dog.CoatColor,
		CoatLength: dog.CoatLength,
		IsLost:     dog.IsLost,
		Latitude:   dog.Latitude,
		Longitude:  dog.Longitude,
		ImgUrl:     dog.ImgUrl,
	}
	if dog.ID != 0 {
		dogModel.ID = dog.ID
	}
	if dog.Owner != nil {
		dogModel.OwnerID = dog.Owner.ID
	}
	if dog.Host != nil {
		dogModel.HostID = dog.Owner.ID
	}
	return dogModel
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
