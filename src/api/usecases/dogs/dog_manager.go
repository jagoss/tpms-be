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

	if exists, _ := d.dogPersister.DogExisitsByNameAndOwner(dog.Name, dog.OwnerID); exists {
		return nil, fmt.Errorf("dog with name %s and ownerID %s already exists", dog.Name, dog.Owner.ID)
	}

	err = setHostAndOwner(dog, userManager)
	if err != nil {
		return nil, err
	}
	log.Printf("[dogmanager.Register] host: %v, owner: %v", dog.Host, dog.Owner)

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

func (d *DogManager) Delete(dogID uint) (bool, error) {
	if err := d.dogPersister.DeleteDog(dogID); err != nil {
		return false, fmt.Errorf("[dogmanager.Delete] error registing dog with id %d: %v", dogID, err)
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
	owner, err := userManager.Get(dog.OwnerID)
	if err != nil {
		return fmt.Errorf("error getting dog owner: %v", err)
	}

	host, err := userManager.Get(dog.HostID)
	if err != nil {
		return fmt.Errorf("error getting dog host: %v", err)
	}
	log.Printf("[dogmanager.setHostAndOwner] host: %v, owner: %v", host, owner)
	dog.Owner = owner
	dog.Host = host
	return nil
}
