package service

import (
	"context"
	"entetry/gotest/internal/dto"
	"entetry/gotest/internal/model"
	"entetry/gotest/internal/repository"
	"github.com/google/uuid"
)

type Company struct {
	companyRepository *repository.Company
}

func NewCompany(companyRepository *repository.Company) *Company {
	return &Company{
		companyRepository: companyRepository}
}

func (c *Company) GetById(ctx context.Context, uuid uuid.UUID) (*model.Company, error) {
	return c.companyRepository.Get(ctx, uuid)
}

func (c *Company) Create(ctx context.Context, request *dto.AddCompanyRequest) error {
	company := &model.Company{ID: uuid.New(), Name: request.Name}
	return c.companyRepository.Create(ctx, company)
}

func (c *Company) Update(ctx context.Context, request *dto.UpdateCompanyRequest) error {
	company := &model.Company{ID: request.Uuid, Name: request.Name}
	return c.companyRepository.Update(ctx, company)
}

func (c *Company) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.companyRepository.Delete(ctx, uuid)
}
