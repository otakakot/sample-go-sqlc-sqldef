// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0

package schema

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Todo struct {
	ID        uuid.UUID          `json:"id"`
	UserID    uuid.UUID          `json:"user_id"`
	Title     string             `json:"title"`
	Completed bool               `json:"completed"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}

type User struct {
	ID        uuid.UUID          `json:"id"`
	Name      string             `json:"name"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}