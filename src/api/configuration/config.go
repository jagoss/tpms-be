package configuration

import (
	"log"
	"os"
)

const (
	DefaultPort = "8080"
	Prod        = "prod"
	Dev         = "dev"
	Test        = "test"
)

type GeneralConfiguration struct {
	Database   DBConfig
	CVModelURL string
}

type DBConfig struct {
	Driver   string
	Host     string
	Username string
	Password string
	Database string
	Port     string
	SSLMode  string
}

func LoadConfiguration() GeneralConfiguration {
	log.Println(os.Getenv("DB_HOST"))
	return GeneralConfiguration{
		Database: DBConfig{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST"),
			Username: os.Getenv("DB_USR"),
			Password: os.Getenv("DB_PWD"),
			Database: os.Getenv("DB_NAME"),
			Port:     os.Getenv("DB_PORT"),
			SSLMode:  os.Getenv("SSLMODE"),
		},
		CVModelURL: os.Getenv("CVMODEL_URL"),
	}
}

func LoadTestConfiguration() GeneralConfiguration {
	return GeneralConfiguration{
		Database: DBConfig{
			Driver:   os.Getenv("DB_DRIVER"),
			Host:     os.Getenv("DB_HOST_TEST"),
			Username: os.Getenv("DB_USR_TEST"),
			Password: os.Getenv("DB_PWD_TEST"),
			Database: os.Getenv("DB_NAME_TEST"),
			Port:     os.Getenv("DB_PORT_TEST"),
			SSLMode:  os.Getenv("SSLMODE"),
		},
		CVModelURL: os.Getenv("CVMODEL_URL"),
	}
}
