package handlers

import "github.com/google/uuid"

type AddCompanyRequest struct {
	Name string `json:"name" validate:"required"`
}

type UpdateCompanyRequest struct {
	UUID uuid.UUID `json:"uuid" validate:"required"`
	Name string    `json:"name" validate:"required"`
}

type AddLogoRequest struct {
	companyID uuid.UUID `json:"uuid" validate:"required"`
	picture   []byte    `json:"picture" validate:"required"`
}

type GetCompanyLogoRequest struct {
	companyID uuid.UUID `json:"uuid" validate:"required"`
}

type CompanyLogoResponse struct {
	picture []byte `json:"picture"`
}
