package mongodb

import (
	"context"
	"entetry/gotest/internal/model"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Company struct {
	db *mongo.Collection
}

func NewCompanyRepository(db *mongo.Database) *Company {
	return &Company{
		db: db.Collection("company"),
	}
}

func (c *Company) GetAll(ctx context.Context) ([]*model.Company, error) {
	cursor, err := c.db.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var result []*model.Company

	for cursor.Next(ctx) {
		company := new(model.Company)
		if err := cursor.Decode(company); err != nil {
			return nil, err
		}
		result = append(result, company)
	}
	err = cursor.Close(ctx)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Company) GetOne(ctx context.Context, uuid uuid.UUID) (*model.Company, error) {
	company := &model.Company{}
	err := c.db.FindOne(ctx, bson.M{"_id": uuid}).Decode(company)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, echo.ErrNotFound
	} else if err != nil {
		return nil, err
	}
	return company, nil
}

func (c *Company) Create(ctx context.Context, company *model.Company) (uuid.UUID, error) {
	company.ID = uuid.New()
	_, err := c.db.InsertOne(ctx, company)
	if err != nil {
		return company.ID, fmt.Errorf("cannot create Company: %v", err)
	}
	return company.ID, err
}

func (c *Company) Update(ctx context.Context, company *model.Company) error {
	r, err := c.db.UpdateOne(ctx, bson.M{"_id": company.ID}, bson.M{"$set": bson.M{"name": company.Name}})
	if err != nil {
		return fmt.Errorf("cannot update Company: %v", err)
	}
	if r.MatchedCount == 0 {
		return echo.ErrNotFound
	}
	return nil
}

func (c *Company) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := c.db.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("cannot delete Company: %v", err)
	}
	return nil
}
