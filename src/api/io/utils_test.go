package io

import (
	"log"
	"testing"
)

func TestToArray(t *testing.T) {
	input := map[string]interface{}{
		"dogId":        23,
		"possibleDogs": "[44]",
	}

	result := ToArray(input["possibleDogs"])
	log.Printf("%s", result)
	if len(result) != 1 && result[0] != "44" {
		t.Fail()
	}
}
