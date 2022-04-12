package dto

import "parcel-service/internal/core/domain"

type CreateDimension struct {
	Width  int
	Height int
	Depth  int
}

type BodyCreateParcel struct {
	OwnerId     string
	Name        string
	Description string
	Size        CreateDimension
	Weight      int
}

type ResponseCreateParcel domain.Parcel

func BuildResponseCreateParcel(model domain.Parcel) ResponseCreateParcel {
	return ResponseCreateParcel(model)
}
