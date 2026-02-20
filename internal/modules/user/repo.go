package user

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepository struct {
	db *pgxpool.Pool
}

var ErrNotFound = errors.New("user not found")

func NewUserRepository(db *pgxpool.Pool) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, u *User) (string, error) {

	query := `INSERT INTO users (id, username, email, password) VALUES ($1, $2, $3, $4) RETURNING id`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var returnedID string
	err := r.db.QueryRow(opCtx, query, u.ID, u.Username, u.Email, u.Password).Scan(&returnedID)
	if err != nil {
		return "", err
	}
	return returnedID, nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (string, error) {
	query := `SELECT id FROM users WHERE email = $1`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var returnedID string
	err := r.db.QueryRow(opCtx, query, email).Scan(&returnedID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", ErrNotFound
		}
		return "", err
	}
	return returnedID, nil
}
