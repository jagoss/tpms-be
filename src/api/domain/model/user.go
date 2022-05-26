package model

type User struct {
	Id          string
	FirstName   string
	LastName    string
	Email       string
	Phone       string
	City        string
	OwnedDogs   []Dog
	HoldingDogs []Dog
}
