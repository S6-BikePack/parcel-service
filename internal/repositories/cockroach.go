package repositories

import (
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"parcel-service/internal/core/domain"
)

type cockroachdb struct {
	Connection *gorm.DB
}

func NewCockroachDB(db *gorm.DB) (*cockroachdb, error) {
	err := db.AutoMigrate(&domain.Customer{}, &domain.Parcel{})

	if err != nil {
		return nil, err
	}

	database := cockroachdb{
		Connection: db,
	}

	return &database, nil
}

func (repository *cockroachdb) GetAll() ([]domain.Parcel, error) {
	var parcels []domain.Parcel

	repository.Connection.Find(&parcels)

	return parcels, nil
}

func (repository *cockroachdb) Get(id string) (domain.Parcel, error) {
	var parcel domain.Parcel

	repository.Connection.Preload(clause.Associations).First(&parcel, "id = ?", id)

	if (parcel == domain.Parcel{}) {
		return parcel, errors.New("parcel not found")
	}

	return parcel, nil
}

func (repository *cockroachdb) GetAllFromCustomer(customerId string) ([]domain.Parcel, error) {
	var parcels []domain.Parcel

	repository.Connection.Where("owner_Id = ?", customerId).Find(&parcels)

	return parcels, nil
}

func (repository *cockroachdb) Save(parcel domain.Parcel) (domain.Parcel, error) {
	result := repository.Connection.Create(&parcel)

	if result.Error != nil {
		return domain.Parcel{}, result.Error
	}

	return parcel, nil
}

func (repository *cockroachdb) Update(parcel domain.Parcel) (domain.Parcel, error) {
	result := repository.Connection.Model(&parcel).Updates(parcel)

	if result.Error != nil {
		return domain.Parcel{}, result.Error
	}

	return parcel, nil
}

func (repository *cockroachdb) SaveOrUpdateCustomer(customer domain.Customer) error {
	updateResult := repository.Connection.Model(&customer).Where("id = ?", customer.ID).Updates(&customer)

	if updateResult.RowsAffected == 0 {
		createResult := repository.Connection.Create(&customer)

		if createResult.Error != nil {
			return errors.New("could not create customer")
		}
	}

	if updateResult.Error != nil {
		return errors.New("could not update customer")
	}

	return nil
}

func (repository *cockroachdb) GetCustomer(id string) (domain.Customer, error) {
	var customer domain.Customer

	repository.Connection.Preload(clause.Associations).First(&customer, "id = ?", id)

	if customer.ID == "" {
		return customer, errors.New("could not find customer with id: " + id)
	}

	return customer, nil
}
