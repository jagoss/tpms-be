package db

import (
	"be-tpms/src/api/configuration"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
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
	log.Printf("host: %s, port: %s, database: %s", host, port, database)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", username, password, host, port, database)

	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})

	if err != nil {
		err = errors.New("Error opening connection to " + host + " database " + database + ". Error: '" + err.Error())
		log.Printf("%v", err)
		return nil, err
	}

	return &DataBase{db}, nil
}
