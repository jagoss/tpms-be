package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"database/sql"
	"log"
	"strings"
)

type PossibleMatchPersister struct {
	connection *db.Connection
}

func NewPossibleMatchPersister(connection *db.Connection) *PossibleMatchPersister {
	return &PossibleMatchPersister{connection}
}

func (pmp *PossibleMatchPersister) AddPossibleMatch(dogID uint, possibleDogID uint) error {
	query := "INSERT INTO tpms_prod.possible_matches(dog_id, possible_dog_id, ack) VALUES (?, ?, 0)"
	_, err := pmp.connection.DB.Exec(query, dogID, possibleDogID)
	if err != nil {
		log.Printf("[PossibleMatchPersister.AddPossibleMatch] error inserting possible match: %s", err.Error())
		return err
	}
	return nil

}
func (pmp *PossibleMatchPersister) UpdateAck(dogID uint, possibleDogID uint, ack model.Ack) error {
	query := "UPDATE tpms_prod.possible_matches SET ack = ? WHERE dog_id = ? AND possible_dog_id = ?"
	log.Printf("[UpdateAck] stm: %s", query)
	log.Printf("[UpdateAck] query values: %d, %d, %d", ack.Int(), dogID, possibleDogID)
	_, err := pmp.connection.DB.Exec(query, ack.Int(), dogID, possibleDogID)
	if err != nil {
		log.Printf("[UpdateAck] error executing stm: %s", err.Error())
		return err
	}
	return nil
}
func (pmp *PossibleMatchPersister) Delete(dogID uint, possibleDogID uint) error {
	query := "DELETE FROM tpms_prod.possible_matches WHERE dog_id = ? AND possible_dog_id = ?"
	log.Printf("[Delete] query: %s", query)
	_, err := pmp.connection.DB.Exec(query, dogID, possibleDogID)
	if err != nil {
		return err
	}
	return nil
}

func (pmp *PossibleMatchPersister) RemovePossibleMatchesForDog(dogID uint) ([]model.PossibleMatch, error) {
	query := "SELECT * FROM tpms_prod.possible_matches WHERE dog_id = ?"
	rows, err := pmp.connection.DB.Query(query, dogID)
	if err != nil {
		return nil, err
	}
	resultList, err := parsePossibleMatch(rows)
	if err != nil {
		return nil, err
	}

	deleteQuery := "DELETE FROM tpms_prod.possible_matches WHERE dog_id IN (?)"
	_, err = pmp.connection.DB.Exec(deleteQuery, dogID)
	if err != nil {
		return nil, err
	}
	return resultList, nil
}

// RemovePossibleDogMatches Remove entries where given id is the possibleDogID
func (pmp *PossibleMatchPersister) RemovePossibleDogMatches(possibleDogID uint) ([]model.PossibleMatch, error) {
	query := "SELECT * FROM tpms_prod.possible_matches WHERE possible_dog_id = ?"
	rows, err := pmp.connection.DB.Query(query, possibleDogID)
	if err != nil {
		return nil, err
	}
	resultList, err := parsePossibleMatch(rows)
	if err != nil {
		return nil, err
	}

	deleteQuery := "DELETE FROM tpms_prod.possible_matches WHERE possible_dog_id IN (?)"
	_, err = pmp.connection.DB.Exec(deleteQuery, possibleDogID)
	if err != nil {
		return nil, err
	}

	return resultList, nil
}

func (pmp *PossibleMatchPersister) GetPossibleMatches(id uint, acks []model.Ack) ([]model.PossibleMatch, error) {
	values := make([]interface{}, len(acks)+2)
	values[0], values[1] = id, id
	for i, ack := range acks {
		values[i+2] = ack.Int()
	}

	query := "SELECT * FROM tpms_prod.possible_matches WHERE (dog_id = ? OR possible_dog_id = ?) AND ack IN (?" + strings.Repeat(",?", len(acks)-1) + ")"
	log.Printf("query: %s", query)
	log.Printf("args: %v", values)
	rows, err := pmp.connection.DB.Query(query, values...)
	if err != nil {
		return nil, err
	}
	var resultList []model.PossibleMatch
	for rows.Next() {
		var pm model.PossibleMatch
		if err := rows.Scan(&pm.DogID, &pm.PossibleDogID, &pm.Ack); err != nil {
			return nil, err
		}
		resultList = append(resultList, pm)
	}

	log.Printf("result list: %v", resultList)

	if resultList == nil || len(resultList) == 0 {
		return make([]model.PossibleMatch, 0), nil
	}
	if err != nil {
		return nil, err
	}
	return resultList, nil
}

func parsePossibleMatch(rows *sql.Rows) ([]model.PossibleMatch, error) {
	var resultList []model.PossibleMatch

	for rows.Next() {
		var pm model.PossibleMatch
		if err := rows.Scan(&pm.DogID, &pm.PossibleDogID, &pm.Ack); err != nil {
			return nil, err
		}
		resultList = append(resultList, pm)
	}

	if resultList == nil || len(resultList) == 0 {
		return make([]model.PossibleMatch, 0), nil
	}
	return resultList, nil
}
