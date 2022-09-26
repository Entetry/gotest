package request

import "github.com/google/uuid"

type AddCompanyRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCompanyRequest struct {
	Uuid uuid.UUID `json:"uuid" validate:"required"`
	Name string    `json:"name" validate:"required"`
}
