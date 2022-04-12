package ports

import "parcel-service/internal/core/domain"

type ParcelService interface {
	GetAll() ([]domain.Parcel, error)
	Get(id string) (domain.Parcel, error)
	GetAllFromCustomer(customerId string) ([]domain.Parcel, error)
	Create(owner, name, description string, size domain.Dimensions, weight int) (domain.Parcel, error)
	UpdateParcelStatus(id string, status int) (domain.Parcel, error)
	CancelParcel(id string) error
	SaveOrUpdateCustomer(customer domain.Customer) error
}
