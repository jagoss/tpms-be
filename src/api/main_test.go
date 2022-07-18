package main

import (
	"be-tpms/src/api/configuration"
	"testing"
)

func TestInitializeDependenciesOk(t *testing.T) {
	_, err := initializeDependencies("../../" + configuration.ConfigFilePath)
	if err != nil {
		t.Fatalf("test failed. Initializing dependencies: %v", err)
	}
}
