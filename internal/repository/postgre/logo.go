package postgre

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"

	"entetry/gotest/internal/model"
)

type LogoRepository interface {
	Create(ctx context.Context, companyId uuid.UUID, image []byte) error
	GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*model.Logo, error)
}

type Logo struct {
	db *pgxpool.Pool
}

func NewLogoRepository(db *pgxpool.Pool) *Logo {
	return &Logo{
		db: db}
}

func (l *Logo) Create(ctx context.Context, companyId uuid.UUID, image []byte) error {
	var logo model.Logo
	logo.ID = uuid.New()
	logo.CompanyID = companyId
	logo.Image = image
	_, err := l.db.Exec(ctx, `INSERT INTO logo (id, company_id, image) VALUES ($1, $2, $3)`, logo.ID,
		logo.CompanyID, logo.Image)

	if err != nil {
		return fmt.Errorf("cannot create Logo: %v", err)
	}
	return nil
}

func (l *Logo) GetByCompanyID(ctx context.Context, companyID uuid.UUID) (*model.Logo, error) {
	var logo model.Logo
	err := l.db.QueryRow(ctx, "SELECT id, company_id, image FROM logo WHERE company_id = $1", companyID).
		Scan(&logo.ID, &logo.CompanyID, &logo.Image)
	if err == pgx.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("GetByCompanyID failed: %v", err)
	}
	return &logo, nil
}
