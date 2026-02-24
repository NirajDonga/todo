package user

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type userRepo struct {
	db *pgxpool.Pool
}

var ErrNotFound = errors.New("user not found")

type UserRepo interface {
	CreateUser(ctx context.Context, u *User) (string, error)
	FindByEmail(ctx context.Context, email string) (*User, error)
}

func NewUserRepository(db *pgxpool.Pool) UserRepo {
	return &userRepo{db: db}
}

func (r *userRepo) CreateUser(ctx context.Context, u *User) (string, error) {

	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var returnedID string
	err := r.db.QueryRow(opCtx, query, u.Username, u.Email, u.Password).Scan(&returnedID)
	if err != nil {
		return "", err
	}
	return returnedID, nil
}

func (r *userRepo) FindByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, username, email, password, created_at FROM users WHERE email = $1`

	opCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	var u User
	err := r.db.QueryRow(opCtx, query, email).Scan(&u.ID, &u.Username, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &u, nil
}
