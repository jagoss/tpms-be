package model

// swagger:model User
type User struct {
	// User ID
	// in: string
	ID string `gorm:"primarykey"`
	// First name
	// in: string
	FirstName string
	// Last name
	// in: string
	LastName string
	// Email
	// in: string
	Email string
	// Phone number
	// in: string
	Phone string
	// City
	// in: string
	City string
}
