package model

import "database/sql"

type Tensor struct {
	Values [1][224][224][3]uint8
}

type DogVector struct {
	ID     uint
	Vector []float64
}

type DogVectorDto struct {
	ID     uint
	Vector sql.NullString
}
