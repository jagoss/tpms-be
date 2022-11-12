package lostandfound

import (
	"log"
	"testing"
)

func TestDistance(t *testing.T) {
	log.Printf("Distance: %.2fkm", Distance(32, 32, 45, 57))
}
