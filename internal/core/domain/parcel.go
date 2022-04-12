package domain

import "errors"

type Parcel struct {
	ID          string
	OwnerId     string
	Name        string
	Description string
	Size        Dimensions `gorm:"embedded"`
	Weight      int
	Status      int
	ServiceArea int
}

func NewParcel(owner, name, description string, size Dimensions, weight, serviceArea int) (Parcel, error) {
	if name == "" {
		return Parcel{}, errors.New("parcel requires a name")
	}

	if size == (Dimensions{}) {
		return Parcel{}, errors.New("parcels dimensions can not be empty")
	}

	if weight == 0 {
		return Parcel{}, errors.New("parcels weight can not be empty")
	}

	return Parcel{
		OwnerId:     owner,
		Name:        name,
		Description: description,
		Size:        size,
		Weight:      weight,
		Status:      0,
		ServiceArea: serviceArea,
	}, nil
}
