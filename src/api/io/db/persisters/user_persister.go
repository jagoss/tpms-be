package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"fmt"
)

type UserPersister struct {
	db *db.DataBase
}

func NewUserPersister(db *db.DataBase) *UserPersister {
	return &UserPersister{db: db}
}

func (up *UserPersister) InsertUser(user *model.User) (*model.User, error) {
	result := up.db.Connection.Create(user)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, fmt.Errorf("no user was inserted")
	}
	return user, nil
}

func (up *UserPersister) GetUser(userID string) (*model.User, error) {
	user := &model.User{ID: userID}
	tx := up.db.Connection.First(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}

func (up *UserPersister) UpdateUser(user *model.User) (*model.User, error) {
	tx := up.db.Connection.Save(user)
	if tx.Error != nil {
		return nil, tx.Error
	}
	return user, nil
}
func (up *UserPersister) DeleteUser(userID string) error {
	tx := up.db.Connection.Delete(&model.User{}, userID)
	return tx.Error
}
