// Package services contains business logic, independent of HTTP concerns.
package services

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/models"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/repository"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/util"
)

// TodoService implements the business logic for managing todos.
type TodoService struct {
	repo *repository.TodoRepository
}

func NewTodoService(repo *repository.TodoRepository) *TodoService {
	return &TodoService{repo: repo}
}

// Create adds a new todo owned by userID.
func (s *TodoService) Create(ctx context.Context, userID, title string) (*models.Todo, error) {
	todo := &models.Todo{
		UserID: userID,
		Title:  title,
	}
	if err := s.repo.Create(ctx, todo); err != nil {
		return nil, err
	}
	return todo, nil
}

// List returns all todos owned by userID.
func (s *TodoService) List(ctx context.Context, userID string) ([]models.Todo, error) {
	return s.repo.ListByUser(ctx, userID)
}

// ToggleCompleted flips is_completed on the todo idHex, scoped to userID.
func (s *TodoService) ToggleCompleted(ctx context.Context, userID, idHex string) (*models.Todo, error) {
	id, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return nil, util.ErrInvalidID
	}
	return s.repo.ToggleCompleted(ctx, userID, id)
}

// ResetCompleted sets is_completed to false for every todo owned by userID.
func (s *TodoService) ResetCompleted(ctx context.Context, userID string) error {
	return s.repo.ResetCompleted(ctx, userID)
}

// Delete removes the todo idHex, scoped to userID.
func (s *TodoService) Delete(ctx context.Context, userID, idHex string) error {
	id, err := bson.ObjectIDFromHex(idHex)
	if err != nil {
		return util.ErrInvalidID
	}
	return s.repo.Delete(ctx, userID, id)
}
