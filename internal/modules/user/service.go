package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/NirajDonga/todo/internal/auth"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUserService(ctx context.Context, u RegisterUser) (UserResponse, error)
	LoginUserService(ctx context.Context, u LoginUser) (AuthResult, error)
}

type userService struct {
	repo UserRepo
	auth auth.AuthService
}

func NewUserService(repo UserRepo, authSvc auth.AuthService) UserService {
	return &userService{
		repo: repo,
		auth: authSvc,
	}
}

func (svc *userService) RegisterUserService(ctx context.Context, u RegisterUser) (UserResponse, error) {
	if u.Email == "" || u.Username == "" || u.Password == "" {
		return UserResponse{}, errors.New("Missing Required Fields")
	}

	existing, err := svc.repo.FindByEmail(ctx, u.Email)
	if err == nil && existing != nil {
		return UserResponse{}, errors.New("Email already registered")
	}
	if err != nil && err != ErrNotFound {
		return UserResponse{}, err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserResponse{}, fmt.Errorf("hashing password failed: %w", err)
	}

	user := User{
		Email:     u.Email,
		Password:  string(hash),
		Username:  u.Username,
		CreatedAt: time.Now(),
	}

	id, err := svc.repo.CreateUser(ctx, &user)
	if err != nil {
		return UserResponse{}, err
	}

	return UserResponse{ID: id, Email: user.Email, Username: user.Username}, nil
}

func (svc *userService) LoginUserService(ctx context.Context, in LoginUser) (AuthResult, error) {
	if in.Email == "" || in.Password == "" {
		return AuthResult{}, errors.New("missing credentials")
	}

	u, err := svc.repo.FindByEmail(ctx, in.Email)
	if err != nil {
		if err == ErrNotFound {
			return AuthResult{}, errors.New("invalid credentials")
		}
		return AuthResult{}, err
	}

	if bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(in.Password)) != nil {
		return AuthResult{}, errors.New("invalid credentials")
	}

	token, err := svc.auth.GenerateToken(ctx, u.ID.String())
	if err != nil {
		return AuthResult{}, fmt.Errorf("generate token: %w", err)
	}

	res := AuthResult{
		Token: token,
		UserInfo: UserResponse{
			ID:       u.ID.String(),
			Email:    u.Email,
			Username: u.Username,
		},
	}
	return res, nil
}
