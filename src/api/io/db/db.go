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

func Init(config configuration.DBConfig, profile string) (*DataBase, error) {
	host := config.Host
	username := config.Username
	password := config.Password
	database := config.Database
	var connectionString string
	if profile == configuration.Prod {
		connectionString = fmt.Sprintf("h2://%s@%s/%s?mem=true", username, host, database)
	} else {
		connectionString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=true", username, password, host, database)
	}

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		return nil, errors.New("Error opening connection to " + config.Driver + " database. Error: '" + err.Error())
	}

	return &DataBase{db}, nil
}
