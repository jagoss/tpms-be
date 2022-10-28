package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
)

type PossibleMatchPersister struct {
	db *db.DataBase
}

func NewPossibleMatchPersister(db *db.DataBase) *PossibleMatchPersister {
	return &PossibleMatchPersister{db}
}

func (pmp *PossibleMatchPersister) AddPossibleMatch(dogID int, possibleDogID int) error {
	return nil
}
func (pmp *PossibleMatchPersister) UpdateAck(dogID int, possibleDogID int, ack model.Ack) error {
	return nil
}
func (pmp *PossibleMatchPersister) Delete(dogID int, possibleDogID int) error { return nil }

// RemovePossibleMatchesForDog Remove entries where given id is the dogID
func (pmp *PossibleMatchPersister) RemovePossibleMatchesForDog(dogID int) error { return nil }

// RemovePossibleDogMatches Remove entries where given id is the possibleDogID
func (pmp *PossibleMatchPersister) RemovePossibleDogMatches(possibleDogID int) error { return nil }
