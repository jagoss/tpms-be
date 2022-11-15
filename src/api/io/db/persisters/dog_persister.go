package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"database/sql"
	"fmt"
	"log"
	"time"
)

const columns = "id, name, breed, age, size, coat_color, coat_length, tail_length, ear, additional_info, is_lost, owner_id, host_id, latitude, longitude, img_url, created_at, deleted_at"

type DogPersister struct {
	connection *db.Connection
}

func NewDogPersister(connection *db.Connection) *DogPersister {
	return &DogPersister{connection}
}

func (dp *DogPersister) InsertDog(dog *model.Dog) (*model.Dog, error) {
	dogModel := mapToDogModel(*dog)

	query := "INSERT INTO tpms_prod.dogs(name, breed, age, size, coat_color, coat_length, tail_length, ear, is_lost, owner_id, host_id, latitude, longitude, img_url, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	result, err := dp.connection.DB.Exec(query, dogModel.Name, dogModel.Breed, dogModel.Age, dogModel.Size, dogModel.CoatColor, dogModel.CoatLength, dogModel.TailLength, dogModel.Ear, dogModel.IsLost, dogModel.OwnerID, dogModel.HostID, dogModel.Latitude, dogModel.Longitude, dogModel.ImgUrl, time.Now())
	if err != nil {
		return nil, err
	}

	dog.ID, _ = result.LastInsertId()
	return dog, nil
}

func (dp *DogPersister) GetDog(dogID uint) (*model.Dog, error) {
	query := fmt.Sprintf("SELECT %s FROM tpms_prod.dogs WHERE id = ?", columns)
	row := dp.connection.DB.QueryRow(query, dogID)
	if row.Err() != nil {
		return nil, row.Err()
	}
	dog, err := parseToDogModel(row)
	if err != nil {
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

	resultDog := mapToDog(*dog, owner, host)

	return &resultDog, nil
}

func (dp *DogPersister) GetDogs(ids []uint) ([]model.Dog, error) {
	query := fmt.Sprintf("SELECT %s FROM tpms_prod.dogs WHERE id in (?)", columns)
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
	query := "UPDATE tpms_prod.dogs SET `name` = ?, age = ?, breed = ?, size = ?, coat_color=?, coat_length = ?, tail_length = ?, ear = ?, is_lost = ?, latitude = ?, longitude = ?, img_url = ? WHERE id = ?"
	_, err := dp.connection.DB.Exec(query, dogModel.Name, dogModel.Age, dogModel.Breed, dogModel.Size, dogModel.CoatColor, dogModel.CoatLength, dogModel.TailLength, dog.Ear, dogModel.IsLost, dogModel.Latitude, dogModel.Longitude, dogModel.ImgUrl, dogModel.ID)
	if err != nil {
		return nil, err
	}

	return dog, nil
}

func (dp *DogPersister) DeleteDog(dogID uint) error {
	query := "DELETE FROM tpms_prod.dogs WHERE id = ?"
	_, err := dp.connection.DB.Exec(query, dogID)
	if err != nil {
		return err
	}

	return nil
}

func (dp *DogPersister) DogExisitsByNameAndOwner(dogName string, ownerID string) (bool, error) {
	query := fmt.Sprintf("SELECT %s FROM tpms_prod.dogs WHERE name = ? AND owner_id = ?", columns)
	rows, err := dp.connection.DB.Query(query, dogName, ownerID)
	if err != nil {
		return false, err
	}
	return rows.Next(), nil
}

func (dp *DogPersister) GetMissingDogs() ([]model.Dog, error) {
	query := fmt.Sprintf("SELECT %s FROM tpms_prod.dogs where is_lost = true", columns)
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
	query := fmt.Sprintf("SELECT DISTINCT %s FROM tpms_prod.dogs WHERE host_id = ? OR owner_id = ?", columns)

	rows, err := dp.connection.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}

	dogs, err := dp.parseDogs(rows)
	if err != nil {
		return nil, err
	}
	return dogs, nil
}

func (dp *DogPersister) SetLostDog(id uint, lat float64, lng float64) error {
	query := "UPDATE tpms_prod.dogs SET latitude = ?, longitude = ? WHERE id = ?"
	_, err := dp.connection.DB.Exec(query, lat, lng, id)
	if err != nil {
		return err
	}
	return nil
}

func (dp *DogPersister) UpdateEmbedding(dogID uint, embedding string) error {
	query := "UPDATE tpms_prod.dogs SET embedding = ? WHERE id = ?"
	_, err := dp.connection.DB.Exec(query, embedding, dogID)
	if err != nil {
		return err
	}
	return nil
}

func (dp *DogPersister) GetPossibleMatchingDog(dog *model.Dog) ([]model.DogVector, error) {
	query := ""
	if dog.Owner != nil {
		query = "SELECT id, embedding FROM tpms_prod.dogs WHERE id != ? AND is_lost = TRUE AND embedding IS NOT NULL AND host_id != ''"
	}
	if dog.Host != nil {
		query = "SELECT id, embedding FROM tpms_prod.dogs WHERE id != ? AND is_lost = TRUE AND embedding IS NOT NULL AND owner_id != ''"
	}
	log.Printf("Query: %s", query)
	rows, err := dp.connection.DB.Query(query, dog.ID)
	if err != nil {
		return nil, err
	}

	var result []model.DogVector
	for rows.Next() {
		var dogVectorDto model.DogVectorDto
		if err := rows.Scan(&dogVectorDto.ID, &dogVectorDto.Vector); err != nil {
			return nil, err
		}
		if !dogVectorDto.Vector.Valid {
			return nil, fmt.Errorf("nil embedding for dog %d", dogVectorDto.ID)
		}
		result = append(result, model.DogVector{ID: dogVectorDto.ID, Vector: ToFloat64List(dogVectorDto.Vector.String)})
	}
	return nil, nil
}

func mapToDog(dogModel model.DogModel, owner *model.User, host *model.User) model.Dog {
	return model.Dog{ID: dogModel.ID,
		Name:           dogModel.Name,
		Breed:          dogModel.Breed,
		Size:           dogModel.Size,
		Age:            dogModel.Age,
		CoatColor:      dogModel.CoatColor,
		CoatLength:     dogModel.CoatLength,
		TailLength:     dogModel.TailLength,
		Ear:            dogModel.Ear,
		IsLost:         dogModel.IsLost,
		AdditionalInfo: dogModel.AdditionalInfo,
		Owner:          owner,
		Host:           host,
		Latitude:       dogModel.Latitude,
		Longitude:      dogModel.Longitude,
		ImgUrl:         dogModel.ImgUrl,
		CreateAt:       dogModel.CreateAt,
	}
}

func mapToDogModel(dog model.Dog) model.DogModel {
	dogModel := model.DogModel{
		Name:           dog.Name,
		Breed:          dog.Breed,
		Size:           dog.Size,
		Age:            dog.Age,
		CoatColor:      dog.CoatColor,
		CoatLength:     dog.CoatLength,
		TailLength:     dog.TailLength,
		Ear:            dog.Ear,
		AdditionalInfo: dog.AdditionalInfo,
		IsLost:         dog.IsLost,
		Latitude:       dog.Latitude,
		Longitude:      dog.Longitude,
		ImgUrl:         dog.ImgUrl,
		CreateAt:       dog.CreateAt,
	}
	if dog.ID != 0 {
		dogModel.ID = dog.ID
	}
	if dog.Owner != nil {
		dogModel.OwnerID = dog.Owner.ID
	} else {
		dogModel.OwnerID = ""
	}
	if dog.Host != nil {
		dogModel.HostID = dog.Host.ID
	} else {
		dogModel.HostID = ""
	}
	return dogModel
}

func (dp *DogPersister) parseDogs(rows *sql.Rows) ([]model.Dog, error) {
	var resultList []model.Dog
	up := UserPersister{dp.connection}

	for rows.Next() {
		var dog model.DogModel
		var deleteDate sql.NullTime
		var additionalInfo sql.NullString
		if err := rows.Scan(&dog.ID, &dog.Name, &dog.Breed, &dog.Age, &dog.Size, &dog.CoatColor, &dog.CoatLength, &dog.TailLength, &dog.Ear, &additionalInfo, &dog.IsLost, &dog.OwnerID, &dog.HostID, &dog.Latitude, &dog.Longitude, &dog.ImgUrl, &dog.CreateAt, &deleteDate); err != nil {
			return nil, err
		}
		if additionalInfo.Valid {
			dog.AdditionalInfo = additionalInfo.String
		}
		if deleteDate.Valid {
			dog.DeleteAt = deleteDate.Time
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

func parseToDogModel(row *sql.Row) (*model.DogModel, error) {
	var dog model.DogModel
	var deleteDate sql.NullTime
	var additionalInfo sql.NullString
	if err := row.Scan(&dog.ID, &dog.Name, &dog.Breed, &dog.Age, &dog.Size, &dog.CoatColor, &dog.CoatLength, &dog.TailLength, &dog.Ear, &additionalInfo, &dog.IsLost, &dog.OwnerID, &dog.HostID, &dog.Latitude, &dog.Longitude, &dog.ImgUrl, &dog.CreateAt, &deleteDate); err != nil {
		return nil, err
	}
	if additionalInfo.Valid {
		dog.AdditionalInfo = additionalInfo.String
	}
	if deleteDate.Valid {
		dog.DeleteAt = deleteDate.Time
	}

	return &dog, nil
}
