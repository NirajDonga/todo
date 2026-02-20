package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo *userRepository
}

func NewUserService(repo *userRepository) *userService {
	return &userService{
		repo: repo,
	}
}

func (svc *userService) RegisterUserService(ctx context.Context, u RegisterUser) (UserResponse, error) {
	if u.Email == "" || u.Username == "" || u.Password == "" {
		return UserResponse{}, errors.New("Missing Required Fields")
	}

	_, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == nil {
		return UserResponse{}, errors.New("Email already registered")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("Hashing of password failed: %w", err)
	}

	user := User{
		ID:        uuid.New().String(),
		Email:     u.Email,
		Password:  hash,
		Username:  u.Username,
		CreatedAt: time.Now(),
	}

	created, err := svc.repo.CreateUser(ctx, &user)
	if err != nil {
		return UserResponse{}, err
	}

	return created
}
