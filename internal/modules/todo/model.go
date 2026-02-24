package todo

import "github.com/google/uuid"

type Todo struct {
	ID        uuid.UUID `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Title     string    `json:"title" db:"title"`
	Completed string    `json:"completed" db:"completed"`
}
