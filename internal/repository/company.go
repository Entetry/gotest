package repository

import (
	"context"
	"entetry/gotest/internal/model"
	"github.com/google/uuid"
)

type CompanyRepository interface {
	Create(ctx context.Context, company *model.Company) (uuid.UUID, error)
	Update(ctx context.Context, company *model.Company) error
	Delete(ctx context.Context, uuid uuid.UUID) error
	GetOne(ctx context.Context, uuid uuid.UUID) (*model.Company, error)
	GetAll(ctx context.Context) ([]*model.Company, error)
}
