package persisters

import (
	"be-tpms/src/api/domain/model"
	"be-tpms/src/api/io/db"
	"fmt"
)

type UserPersister struct {
	connection *db.Connection
}

func NewUserPersister(connection *db.Connection) *UserPersister {
	return &UserPersister{connection: connection}
}

func (up *UserPersister) InsertUser(user *model.User) (*model.User, error) {
	query := "INSERT INTO users(`id`, `name`, `phone`, `email`, `latitude`, `longitude`) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := up.connection.DB.Exec(query, user.ID, user.Name, user.Phone, user.Email, user.Latitude, user.Longitude)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (up *UserPersister) GetUser(userID string) (*model.User, error) {
	query := "SELECT * FROM users WHERE id = ?"
	rows, err := up.connection.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}

	var user model.User
	if err = rows.Scan(&user.ID, &user.Email, &user.Phone, &user.FCMToken, &user.Name, &user.Latitude, &user.Longitude); err != nil {
		return nil, err
	}
	return &user, nil
}

func (up *UserPersister) UpdateUser(user *model.User) (*model.User, error) {
	query := "UPDATE users SET email = ?, phone = ?, fcm_token = ?, latitude = ?, longitude = ? WHERE id = ?"
	_, err := up.connection.DB.Exec(query, user.Email, user.Phone, user.FCMToken, user.Latitude, user.Longitude, user.ID)
	if err != nil {
		return nil, err
	}

	return user, nil
}
func (up *UserPersister) DeleteUser(userID string) error {
	query := "DELETE FROM users WHERE id = ?"
	result, err := up.connection.DB.Exec(query, userID)
	if err != nil {
		return err
	}

	if amount, err := result.RowsAffected(); err != nil || amount == 0 {
		if err != nil {
			return err
		}
		return fmt.Errorf("no rows affected in database when deleting user %s", userID)
	}

	return nil
}
