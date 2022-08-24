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
	env, err := initializeDependencies()
	if err != nil {
		panic(any(fmt.Errorf("error Init Main: %w", err)))
	}
	InitRouter(env, "8080")
}

func InitRouter(env *environment.Env, port string) {
	err := router.Run(*env, port)
	if err != nil {
		panic(any(fmt.Errorf("error Running router: %w", err)))
	}
}

func initializeDependencies() (*environment.Env, error) {
	scope := os.Getenv("SCOPE")

	var conf configuration.GeneralConfiguration
	if scope == "" || scope == configuration.Test {
		conf = configuration.LoadTestConfiguration()
	} else {
		conf = configuration.LoadConfiguration()
	}

	firebaseAuth := *configuration.SetupFirebase()
	database, err := initializeDatabase(conf)
	if err != nil {
		return nil, err
	}
	userPersister := persisters.NewUserPersister(database)
	dogPersister := persisters.NewDogPersister(database)
	restClient := *router.CreateRestClientConfig(scope)
	cvModelClient := restclient.NewCVModelRestClient(&restClient)
	env := environment.InitEnv(firebaseAuth, restClient, cvModelClient, userPersister, dogPersister)
	return env, nil
}

func initializeDatabase(config configuration.GeneralConfiguration) (*db.DataBase, error) {
	database, err := db.Init(config.Database)
	if err != nil {
		return nil, fmt.Errorf("unable to init database configuration")
	}
	return database, nil
}
