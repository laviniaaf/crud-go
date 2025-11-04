// models = data structures
package main

import (
	"time"

	"github.com/google/uuid"
)

type Bill struct {
	ID        uuid.UUID `json:"id"`
	Embasa    float64   `json:"embasa"`
	Coelba    float64   `json:"coelba"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
