package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"gorm.io/gorm/clause"
)

type PossibleMatchPersister struct {
	db *db.DataBase
}

func NewPossibleMatchPersister(db *db.DataBase) *PossibleMatchPersister {
	return &PossibleMatchPersister{db}
}

func (pmp *PossibleMatchPersister) AddPossibleMatch(dogID uint, possibleDogID uint) error {
	tx := pmp.db.Connection.Create(model.PossibleMatch{DogID: dogID, PossibleDogID: possibleDogID, Ack: model.Pending})
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (pmp *PossibleMatchPersister) UpdateAck(dogID uint, possibleDogID uint, ack model.Ack) error {
	possibleMatch := &model.PossibleMatch{DogID: dogID, PossibleDogID: possibleDogID}
	tx := pmp.db.Connection.First(possibleMatch)
	if tx.Error != nil {
		return tx.Error
	}

	possibleMatch.Ack = ack
	tx = pmp.db.Connection.Save(possibleMatch)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}
func (pmp *PossibleMatchPersister) Delete(dogID uint, possibleDogID uint) error {
	registerToDelete := &model.PossibleMatch{DogID: dogID, PossibleDogID: possibleDogID}
	tx := pmp.db.Connection.Delete(registerToDelete)
	if tx.Error != nil {
		return tx.Error
	}
	return nil
}

// RemovePossibleMatchesForDog Remove entries where given id is the dogID
func (pmp *PossibleMatchPersister) RemovePossibleMatchesForDog(dogID uint) ([]model.PossibleMatch, error) {
	var possibleMatches []model.PossibleMatch
	tx := pmp.db.Connection.
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "possible_dog_id"}}}).
		Where("dog_id = ?", dogID).
		Delete(possibleMatches)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return possibleMatches, nil
}

// RemovePossibleDogMatches Remove entries where given id is the possibleDogID
func (pmp *PossibleMatchPersister) RemovePossibleDogMatches(possibleDogID uint) ([]model.PossibleMatch, error) {
	var possibleMatches []model.PossibleMatch
	tx := pmp.db.Connection.
		Clauses(clause.Returning{Columns: []clause.Column{{Name: "dog_id"}}}).
		Where("possible_dog_id = ?", possibleDogID).
		Delete(possibleMatches)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return possibleMatches, nil
}
