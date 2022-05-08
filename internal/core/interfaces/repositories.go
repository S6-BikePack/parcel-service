package interfaces

import "parcel-service/internal/core/domain"

type ParcelRepository interface {
	GetAll() ([]domain.Parcel, error)
	Get(id string) (domain.Parcel, error)
	GetAllFromCustomer(customerId string) ([]domain.Parcel, error)
	Save(parcel domain.Parcel) (domain.Parcel, error)
	Update(parcel domain.Parcel) (domain.Parcel, error)
	SaveOrUpdateCustomer(customer domain.Customer) error
	GetCustomer(id string) (domain.Customer, error)
}
