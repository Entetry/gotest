package model

import "github.com/google/uuid"

type Company struct {
	ID   uuid.UUID
	Name string
}
