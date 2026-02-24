package todo

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

type TodoService interface {
	CreateTodoService(ctx context.Context, userID string, title string) (Todo, error)
	GetTodosService(ctx context.Context, userID string) ([]Todo, error)
}

type todoService struct {
	repo TodoRepo
}

func NewTodoService(repo TodoRepo) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) CreateTodoService(ctx context.Context, userID string, title string) (Todo, error) {
	if title == "" {
		return Todo{}, errors.New("title is required")
	}
	uid, err := uuid.Parse(userID)
	if err != nil {
		return Todo{}, err
	}
	id, err := s.repo.CreateTodo(ctx, uid, title)
	if err != nil {
		return Todo{}, err
	}
	t := Todo{
		ID:        id,
		UserID:    uid,
		Title:     title,
		Completed: "false",
	}
	return t, nil
}

func (s *todoService) GetTodosService(ctx context.Context, userID string) ([]Todo, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	return s.repo.GetTodos(ctx, uid)
}
