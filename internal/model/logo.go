package model

import "github.com/google/uuid"

type Logo struct {
	ID        uuid.UUID
	CompanyID uuid.UUID
	Image     []byte
}
