package postgre

import (
	"context"
	"entetry/gotest/internal/model"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Company struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *Company {
	return &Company{
		db: db,
	}
}

func (c *Company) GetAll(ctx context.Context) ([]*model.Company, error) {
	rows, err := c.db.Query(ctx, `SELECT id, name FROM company`)
	if err != nil {
		return nil, fmt.Errorf("query: %v", err)
	}
	defer rows.Close()

	var results []*model.Company

	for rows.Next() {
		var company model.Company

		err = rows.Scan(&company.ID, &company.Name)
		if err != nil {
			return nil, fmt.Errorf("scan: %v", err)
		}

		results = append(results, &company)
	}

	return results, nil
}

func (c *Company) GetOne(ctx context.Context, uuid uuid.UUID) (*model.Company, error) {
	var company *model.Company
	err := c.db.QueryRow(ctx, "SELECT id, name FROM company WHERE id = $1", uuid).Scan(&company.ID, &company.Name)
	return company, err
}

func (c *Company) Create(ctx context.Context, company *model.Company) (uuid.UUID, error) {
	company.ID = uuid.New()
	err := c.db.QueryRow(ctx, "INSERT INTO company(id, name) VALUES ($1, $2) RETURNING id, name;",
		company.ID, company.Name).Scan(&company.ID, &company.Name)
	if err != nil {
		return company.ID, fmt.Errorf("cannot create Company: %v", err)
	}
	return company.ID, err
}

func (c *Company) Update(ctx context.Context, company *model.Company) error {
	_, err := c.db.Exec(ctx, "UPDATE company SET name = $2 WHERE id=$1 RETURNING id, name;",
		company.ID, company.Name)
	if err != nil {
		return fmt.Errorf("cannot update Company: %v", err)
	}
	return err
}

func (c *Company) Delete(ctx context.Context, uuid uuid.UUID) error {
	_, err := c.db.Exec(ctx, "DELETE FROM company WHERE id = $1", uuid)
	if err != nil {
		return fmt.Errorf("cannot delete Company: %v", err)
	}
	return nil
}
