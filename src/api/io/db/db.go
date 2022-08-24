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
	port := config.Port
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, errors.New("Error opening connection to " + config.Driver + " database. Error: '" + err.Error())
	}

	return &DataBase{db}, nil
}
