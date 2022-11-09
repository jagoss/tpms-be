package db

import (
	"be-tpms/src/api/configuration"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Connection struct {
	DB *sql.DB
}

func Init(config configuration.DBConfig) (*Connection, error) {
	host := config.Host
	username := config.Username
	password := config.Password
	database := config.Database
	port := config.Port
	log.Printf("host: %s, port: %s, database: %s", host, port, database)
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true", username, password, host, port, database)

	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		err = errors.New("Error opening connection to " + host + " database " + database + ". Error: '" + err.Error())
		log.Printf("%v", err)
		return nil, err
	}
	return &Connection{db}, nil
}
