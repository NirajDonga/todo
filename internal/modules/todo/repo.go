package todo

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type todorepo struct {
	db *pgxpool.Pool
}

type TodoRepo interface {
	CreateTodo(ctx context.Context, userID uuid.UUID, title string) (uuid.UUID, error)
	GetTodos(ctx context.Context, userID uuid.UUID) ([]Todo, error)
}

func NewTodoRepo(db *pgxpool.Pool) TodoRepo {
	return &todorepo{db: db}
}

func (r *todorepo) CreateTodo(ctx context.Context, userID uuid.UUID, title string) (uuid.UUID, error) {
	query := `INSERT INTO todo (user_id, title) VALUES ($1, $2) RETURNING id`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var returnedID uuid.UUID
	err := r.db.QueryRow(opCtx, query, userID, title).Scan(&returnedID)
	if err != nil {
		return uuid.Nil, err
	}
	return returnedID, nil
}

func (r *todorepo) GetTodos(ctx context.Context, userID uuid.UUID) ([]Todo, error) {
	query := `SELECT id, user_id, title, completed FROM todo WHERE user_id = $1 ORDER BY id`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	rows, err := r.db.Query(opCtx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []Todo
	for rows.Next() {
		var t Todo
		if err := rows.Scan(&t.ID, &t.UserID, &t.Title, &t.Completed); err != nil {
			return nil, err
		}
		out = append(out, t)
	}
	return out, nil
}
