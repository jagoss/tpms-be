package interfaces

import "be-tpms/src/api/domain/model"

type PossibleMatchPersister interface {
	AddPossibleMatch(uint, uint) error
	UpdateAck(uint, uint, model.Ack) error
	Delete(uint, uint) error
	// RemovePossibleMatchesForDog Remove entries where given id is the dogID
	RemovePossibleMatchesForDog(uint) ([]model.PossibleMatch, error)
	// RemovePossibleDogMatches Remove entries where given id is the possibleDogID
	RemovePossibleDogMatches(uint) ([]model.PossibleMatch, error)
	//GetPossibleMatches Return possible matches given dogID either if its possible dog or dog
	GetPossibleMatches(id uint, acks []model.Ack) ([]model.PossibleMatch, error)
}
