package main

import (
	"be-tpms/src/api/configuration"
	"be-tpms/src/api/environment"
	"be-tpms/src/api/io/db"
	"be-tpms/src/api/io/db/persisters"
	"be-tpms/src/api/io/restclient"
	"be-tpms/src/api/io/router"
	"fmt"
	"os"
)

func main() {
	env, err := initializeDependencies(configuration.ConfigFilePath)
	if err != nil {
		panic(any(fmt.Errorf("error Init Main: %v", err)))
	}
	InitRouter(env, "8080")
}

func InitRouter(env *environment.Env, port string) {
	err := router.Run(*env, port)
	if err != nil {
		panic(any(fmt.Errorf("error Running router: %v", err)))
	}
}

func initializeDependencies(configurationPackagePath string) (*environment.Env, error) {
	scope := os.Getenv("SCOPE")
	if scope == "" {
		scope = "test"
	}
	path := configurationPackagePath + "/" + scope + "_config.yml"
	conf := configuration.GeneralConfiguration{}
	err := conf.LoadConfiguration(path)
	if err != nil {
		return nil, fmt.Errorf("error initializing dependencies: %v", err)
	}

	database, err := initializeDatabase(conf, scope)
	if err != nil {
		return nil, err
	}
	userPersister := persisters.NewUserPersister(database)
	dogPersister := persisters.NewDogPersister(database)
	restClient := *router.CreateRestClientConfig(scope)
	cvModelClient := restclient.NewCVModelRestClient(&restClient)

	return &environment.Env{
		RestClient:        restClient,
		CVModelRestClient: cvModelClient,
		UserPersister:     userPersister,
		DogPersister:      dogPersister,
	}, nil
}

func initializeDatabase(config configuration.GeneralConfiguration, scope string) (*db.DataBase, error) {
	database, err := db.Init(config.Database, scope)
	if err != nil {
		return nil, fmt.Errorf("unable to init database configuration")
	}
	return database, nil
}
