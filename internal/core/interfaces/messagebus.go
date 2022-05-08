package interfaces

import "parcel-service/internal/core/domain"

type MessageBusPublisher interface {
	CreateParcel(parcel domain.Parcel) error
	UpdateParcelStatus(id string, status int) error
}
