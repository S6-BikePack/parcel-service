package domain

type Customer struct {
	ID          string
	ServiceArea int
	Parcels     []Parcel `gorm:"foreignKey:OwnerId"`
}

func NewCustomer(id string, serviceArea int) Customer {
	return Customer{
		ID:          id,
		ServiceArea: serviceArea,
	}
}
