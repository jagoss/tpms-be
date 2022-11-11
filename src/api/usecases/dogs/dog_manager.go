package dogs

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
	"log"
)

type DogManager struct {
	dogPersister interfaces.DogPersister
	DOBucket     interfaces.Storage
}

func NewDogManager(dogPersister interfaces.DogPersister, storage interfaces.Storage) *DogManager {
	return &DogManager{
		dogPersister: dogPersister,
		DOBucket:     storage,
	}
}

func (d *DogManager) Get(dogID uint) (*model.Dog, error) {
	dog, err := d.dogPersister.GetDog(dogID)
	if err != nil {
		return nil, fmt.Errorf("[dogmanager.Get] error getting dog with id %d: %v", dogID, err)
	}

	return dog, nil
}

func (d *DogManager) Register(dog *model.Dog, imgBuffArray []string, userManager interfaces.UserManager) (*model.Dog, error) {
	imgsPath, err := d.AddImgs(dog, imgBuffArray)
	if err != nil {
		return nil, err
	}

	dog.ImgUrl = imgsPath
	if dog.Owner != nil {
		if exists, _ := d.dogPersister.DogExisitsByNameAndOwner(dog.Name, dog.Owner.ID); exists {
			return nil, fmt.Errorf("dog with name %s and ownerID %s already exists", dog.Name, dog.Owner.ID)
		}
	}

	err = setHostAndOwner(dog, userManager)
	if err != nil {
		return nil, err
	}
	log.Printf("[dogmanager.Register] host: %v, owner: %v", dog.Host, dog.Owner)
	dog.IsLost = dog.Owner == nil || dog.Host == nil
	newDog, err := d.dogPersister.InsertDog(dog)
	if err != nil {
		return nil, fmt.Errorf("[dogmanager.Register] error registing dog: %v", err)
	}

	return newDog, nil
}

func (d *DogManager) Modify(dog *model.Dog, imgBuffArray []string, userManager interfaces.UserManager) (*model.Dog, error) {
	imgsPath, err := d.AddImgs(dog, imgBuffArray)
	if err != nil {
		return nil, err
	}
	dog.ImgUrl = dog.ImgUrl + ";" + imgsPath

	err = setHostAndOwner(dog, userManager)
	if err != nil {
		return nil, err
	}

	updatedDog, err := d.dogPersister.UpdateDog(dog)
	if err != nil {
		return nil, fmt.Errorf("[dogmanager.Modify] error modifying dog with id %d: %v", dog.ID, err)
	}

	return updatedDog, nil
}

func (d *DogManager) Delete(dogID uint, persister interfaces.PossibleMatchPersister) (bool, error) {
	if _, err := persister.RemovePossibleDogMatches(dogID); err != nil {
		return false, fmt.Errorf("[dogmanager.Delete] error deleting matches where possible dog has id %d: %v", dogID, err)
	}
	if _, err := persister.RemovePossibleMatchesForDog(dogID); err != nil {
		return false, fmt.Errorf("[dogmanager.Delete] error deleting possible matches for dog with id %d: %v", dogID, err)

	}
	if err := d.dogPersister.DeleteDog(dogID); err != nil {
		return false, fmt.Errorf("[dogmanager.Delete] error deleting dog with id %d: %v", dogID, err)
	}

	return true, nil
}

func (d *DogManager) AddImgs(dog *model.Dog, imgBuffArray []string) (string, error) {
	if len(imgBuffArray) > 0 {
		imgsPath, err := d.DOBucket.SaveImgs(imgBuffArray)
		if err != nil {
			return "", fmt.Errorf("[dogmanager.AddImgs] error saving DogImgs %d: %v", dog.ID, err)
		}
		return imgsPath, nil
	}

	return "", nil
}

func (d *DogManager) GetAllUserDogs(userID string) ([]model.Dog, []model.Dog, error) {
	missingDogs, err := d.dogPersister.GetDogsByUser(userID)
	if err != nil {
		return nil, nil, fmt.Errorf("[dogmanaer.GetAllUserDogs] error getting user %s dogs: %s", userID, err.Error())
	}
	var foundDogs []model.Dog
	var userOwnedDogs []model.Dog
	for _, dog := range missingDogs {
		if dog.Owner == nil {
			foundDogs = append(foundDogs, dog)
		} else {
			userOwnedDogs = append(userOwnedDogs, dog)
		}
	}

	return foundDogs, userOwnedDogs, nil
}

func setHostAndOwner(dog *model.Dog, userManager interfaces.UserManager) error {
	if dog.Owner != nil {
		owner, err := userManager.Get(dog.Owner.ID)
		if err != nil {
			return fmt.Errorf("error getting dog owner: %v", err)
		}
		dog.Owner = owner
	}

	if dog.Host != nil {
		host, err := userManager.Get(dog.Owner.ID)
		if err != nil {
			return fmt.Errorf("error getting dog owner: %v", err)
		}
		dog.Host = host
	}

	return nil
}
