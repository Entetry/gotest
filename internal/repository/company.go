package repository

import (
	"context"
	"entetry/gotest/internal/model"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
	log "github.com/sirupsen/logrus"
)

type Company struct {
	db *pgxpool.Pool
}

func NewCompanyRepository(db *pgxpool.Pool) *Company {
	return &Company{
		db: db,
	}
}

func (c *Company) Get(ctx context.Context, uuid uuid.UUID) (*model.Company, error) {
	company := &model.Company{}
	err := c.db.QueryRow(ctx, "SELECT id, name FROM company WHERE id = $1", uuid).Scan(&company.ID, &company.Name)
	return company, err
}

func (c *Company) Create(ctx context.Context, company *model.Company) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.TODO())
			if err != nil {
				log.Errorf("rollback wasn't: %v", err)
			}
		} else {
			err = tx.Commit(context.TODO())
			if err != nil {
				log.Errorf("rollback wasn't: %v", err)
			}
		}
	}()

	err = tx.QueryRow(ctx, "INSERT INTO company(id, name) VALUES ($1, $2) RETURNING id, name;",
		company.ID, company.Name).Scan(&company.ID, &company.Name)
	if err != nil {
		return err
	}
	return err
}

func (c *Company) Update(ctx context.Context, company *model.Company) error {
	tx, err := c.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = tx.Rollback(context.TODO())
			if err != nil {
				log.Errorf("rollback wasn't: %v", err)
			}
		} else {
			err = tx.Commit(context.TODO())
			if err != nil {
				log.Errorf("rollback wasn't: %v", err)
			}
		}
	}()

	err = tx.QueryRow(ctx, "UPDATE company SET name = $2 WHERE id=$1 RETURNING id, name;",
		company.ID, company.Name).Scan(&company.ID, &company.Name)
	if err != nil {
		return err
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
