package model

// swagger:model User
type User struct {
	// User ID
	// in: string
	ID string `json:"id"`
	// Name
	// in: string
	Name string `json:"name"`
	// Email
	// in: string
	Email string `json:"email"`
	// Phone number
	// in: string
	Phone string `json:"phone"`
	// Latitude
	// in: float64
	Latitude float64 `json:"latitude"`
	// Longitude
	// in: float64
	Longitude float64 `json:"longitude"`
	// FCMToken
	// in: string
	FCMToken string `json:"FCMToken"`
}

// swagger:model UserContactInfo
type UserContactInfo struct {
	// Name
	// in: string
	Name string `json:"name"`
	// Email
	// in: string
	Email string `json:"email"`
	// Phone number
	// in: string
	Phone string `json:"phone"`
}
