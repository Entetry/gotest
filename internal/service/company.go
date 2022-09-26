package service

import (
	"context"
	"entetry/gotest/internal/model"
	"entetry/gotest/internal/repository"
	"github.com/google/uuid"
)

type Company struct {
	companyRepository repository.CompanyRepository
}

func NewCompany(companyRepository repository.CompanyRepository) *Company {
	return &Company{
		companyRepository: companyRepository}
}

func (c *Company) GetAll(ctx context.Context) ([]*model.Company, error) {
	return c.companyRepository.GetAll(ctx)
}

func (c *Company) GetById(ctx context.Context, uuid uuid.UUID) (*model.Company, error) {
	return c.companyRepository.GetOne(ctx, uuid)
}

func (c *Company) Create(ctx context.Context, company *model.Company) (uuid.UUID, error) {
	return c.companyRepository.Create(ctx, company)
}

func (c *Company) Update(ctx context.Context, company *model.Company) error {
	return c.companyRepository.Update(ctx, company)
}

func (c *Company) Delete(ctx context.Context, uuid uuid.UUID) error {
	return c.companyRepository.Delete(ctx, uuid)
}
