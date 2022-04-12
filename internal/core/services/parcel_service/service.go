package parcel_service

import (
	"fmt"
	"parcel-service/internal/core/domain"
	"parcel-service/internal/core/ports"
)

type service struct {
	parcelRepository ports.ParcelRepository
	messagePublisher ports.MessageBusPublisher
}

func New(parcelRepository ports.ParcelRepository, messagePublisher ports.MessageBusPublisher) *service {
	return &service{
		parcelRepository: parcelRepository,
		messagePublisher: messagePublisher,
	}
}

func (srv *service) GetAll() ([]domain.Parcel, error) {
	return srv.parcelRepository.GetAll()
}

func (srv *service) Get(id string) (domain.Parcel, error) {
	return srv.parcelRepository.Get(id)
}

func (srv *service) GetAllFromCustomer(customerId string) ([]domain.Parcel, error) {
	customer, err := srv.parcelRepository.GetCustomer(customerId)

	if err != nil {
		return nil, err
	}

	return srv.parcelRepository.GetAllFromCustomer(customer.ID)
}

func (srv *service) Create(owner, name, description string, size domain.Dimensions, weight int) (domain.Parcel, error) {
	customer, err := srv.parcelRepository.GetCustomer(owner)

	if err != nil {
		return domain.Parcel{}, err
	}

	parcel, err := domain.NewParcel(customer.ID, name, description, size, weight, customer.ServiceArea)

	if err != nil {
		return domain.Parcel{}, err
	}

	parcel, err = srv.parcelRepository.Save(parcel)

	if err != nil {
		return domain.Parcel{}, err
	}

	srv.messagePublisher.CreateParcel(parcel)

	return parcel, nil
}

func (srv *service) UpdateParcelStatus(id string, status int) (domain.Parcel, error) {
	parcel, err := srv.Get(id)

	if err != nil {
		return domain.Parcel{}, err
	}

	parcel.Status = status

	parcel, err = srv.parcelRepository.Update(parcel)

	if err != nil {
		return domain.Parcel{}, err
	}

	srv.messagePublisher.UpdateParcelStatus(parcel.ID, parcel.Status)

	return parcel, nil
}

func (srv *service) CancelParcel(id string) error {
	parcel, err := srv.Get(id)

	if err != nil {
		return err
	}

	parcel.Status = -1
	parcel, err = srv.parcelRepository.Update(parcel)

	if err != nil {
		return err
	}

	srv.messagePublisher.UpdateParcelStatus(parcel.ID, parcel.Status)

	return err
}

func (srv *service) SaveOrUpdateCustomer(customer domain.Customer) error {
	fmt.Println(customer)
	err := srv.parcelRepository.SaveOrUpdateCustomer(customer)

	return err
}
