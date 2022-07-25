package model

type User struct {
	ID        string `gorm:"primarykey"`
	FirstName string
	LastName  string
	Email     string
	Phone     string
	City      string
}
