package dogs

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
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

func (d *DogManager) Register(dog *model.Dog, imgBuffArray [][]byte) (*model.Dog, error) {
	imgsPath, err := d.AddImgs(dog, imgBuffArray)
	if err != nil {
		return nil, err
	}
	dog.ImgUrl = imgsPath

	if d.dogPersister.DogExisitsByNameAndOwner(dog.Name, dog.Owner.ID) {
		return nil, fmt.Errorf("dog with name %s and ownerID %s already exists", dog.Name, dog.Owner.ID)
	}

	newDog, err := d.dogPersister.InsertDog(dog)
	if err != nil {
		return nil, fmt.Errorf("[dogmanager.Register] error registing dog: %v", err)
	}
	return newDog, nil
}

func (d *DogManager) Modify(dog *model.Dog, imgBuffArray [][]byte) (*model.Dog, error) {
	imgsPath, err := d.AddImgs(dog, imgBuffArray)
	if err != nil {
		return nil, err
	}
	dog.ImgUrl = dog.ImgUrl + ";" + imgsPath
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

func (d *DogManager) AddImgs(dog *model.Dog, imgBuffArray [][]byte) (string, error) {
	if len(imgBuffArray) > 0 {
		imgsPath, err := d.DOBucket.SaveImgs(imgBuffArray)
		if err != nil {
			return "", fmt.Errorf("[dogmanager.AddImgs] error saving DogImgs %d: %v", dog.ID, err)
		}
		return imgsPath, nil
	}
	return "", nil
}
