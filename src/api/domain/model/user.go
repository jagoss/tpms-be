package model

// swagger:model User
type User struct {
	// User ID
	// in: string
	ID string `gorm:"primarykey"`
	// Name
	// in: string
	Name string
	// Email
	// in: string
	Email string
	// Phone number
	// in: string
	Phone string
	// Latitude
	// in: float64
	Latitude float64
	// Longitude
	// in: float64
	Longitude float64
	// FCMToken
	// in: string
	FCMToken string
}

// swagger:model UserContactInfo
type UserContactInfo struct {
	// Name
	// in: string
	Name string
	// Email
	// in: string
	Email string
	// Phone number
	// in: string
	Phone string
}
