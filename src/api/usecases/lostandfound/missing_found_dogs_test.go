package lostandfound

import (
	"log"
	"testing"
)

func TestDistance(t *testing.T) {
	log.Printf("distance: %.2fkm", distance(32, 32, 45, 57))
}
