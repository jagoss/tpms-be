package model

type User struct {
	ID        string `gorm:"primarykey"`
	FirstName string
	LastName  string
	Email     string
	Phone     string
	City      string
	OwnedDogs []Dog `gorm:"foreignKey:OwnerID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	HostDogs  []Dog `gorm:"foreignKey:HostID,constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
