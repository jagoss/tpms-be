package dogs

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/storage"
	"be-tpms/src/api/usecases/interfaces"
	"fmt"
	"log"
)

type DogManager struct {
	dogPersister interfaces.DogPersister
}

func NewDogManager(dogPersister interfaces.DogPersister) *DogManager {
	return &DogManager{dogPersister: dogPersister}
}

func (d *DogManager) Get(dogID uint) (*model.Dog, error) {
	dog, err := d.dogPersister.GetDog(dogID)
	if err != nil {
		return nil, fmt.Errorf("[dogmanager.Get] error getting dog with id %d: %v", dogID, err)
	}
	return dog, nil
}

func (d *DogManager) Register(dog *model.Dog) (*model.Dog, error) {
	if d.dogPersister.DogExisitsByNameAndOwner(dog.Name, dog.Owner.ID) {
		return nil, fmt.Errorf("dog with name %s and ownerID %s already exists", dog.Name, dog.Owner.ID)
	}
	dog, err := d.dogPersister.InsertDog(dog)
	if err != nil {
		return nil, fmt.Errorf("[dogmanager.Register] error registing dog: %v", err)
	}
	return dog, nil
}

func (d *DogManager) Modify(dog *model.Dog) (*model.Dog, error) {
	dog, err := d.dogPersister.UpdateDog(dog)
	if err != nil {
		return nil, fmt.Errorf("[dogmanager.Modify] error modifying dog with id %d: %v", dog.ID, err)
	}
	return dog, nil
}

func (d *DogManager) Delete(dogID uint) (bool, error) {
	if err := d.dogPersister.DeleteDog(dogID); err != nil {
		return false, fmt.Errorf("[dogmanager.Delete] error registing dog with id %d: %v", dogID, err)
	}
	return true, nil
}

func (d *DogManager) AddImg(imgBuff []byte) error {
	bucket := storage.NewBucket()
	imgName, err := bucket.SaveImgs([][]byte{imgBuff})
	if err != nil {
		return err
	}
	log.Printf("img name: %s", imgName)
	return nil
}
