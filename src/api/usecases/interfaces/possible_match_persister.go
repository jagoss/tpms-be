package interfaces

import "be-tpms/src/api/domain/model"

type PossibleMatchPersister interface {
	AddPossibleMatch(int, int) error
	UpdateAck(int, int, model.Ack) error
	Delete(int, int) error
	RemovePossibleMatchesForDog(int) error
	RemovePossibleDogMatches(int) error
}
