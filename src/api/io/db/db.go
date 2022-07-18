package db

import (
	"be-tpms/src/api/configuration"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DataBase struct {
	Connection *gorm.DB
}

func Init(config configuration.DBConfig) (*DataBase, error) {
	host := config.Host
	username := config.Username
	password := config.Password
	database := config.Database
	connectionString := fmt.Sprintf("%mocks:%mocks@tcp(%mocks)/%mocks?charset=utf8&parseTime=true", username, password, host, database)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, errors.New("Error opening connection to " + config.Driver + " database. Error: '" + err.Error())
	}

	return &DataBase{db}, nil
}
