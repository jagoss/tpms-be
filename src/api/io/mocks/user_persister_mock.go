package mock

import "be-tpms/src/api/domain/model"

type UserPersisterMock struct {
	users         []model.User
	userCounter   int
	errors        []error
	errorsCounter int
}

func (up *UserPersisterMock) InsertUser(user *model.User) (*model.User, error) {
	u := up.users[up.userCounter]
	up.userCounter += 1
	return &u, nil
}

func (up *UserPersisterMock) GetUser(userID string) (*model.User, error) {
	user := up.users[up.userCounter]
	up.userCounter += 1
	return &user, nil
}

func (up *UserPersisterMock) UpdateUser(user *model.User) (*model.User, error) {
	return user, nil
}

func (up *UserPersisterMock) DeleteUser(userID string) error {
	err := up.errors[up.errorsCounter]
	up.errorsCounter += 1
	return err
}
