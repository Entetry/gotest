package model

import "github.com/google/uuid"

type Company struct {
	ID   uuid.UUID `bson:"_id"`
	Name string    `bson:"name"`
}
