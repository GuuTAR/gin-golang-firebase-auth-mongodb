// Package repository contains the database data access layer.
package repository

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/models"
	"github.com/GuuTAR/gin-golang-firebase-auth-mongodb/internal/util"
)

// TodoRepository provides MongoDB-backed CRUD access to todos.
type TodoRepository struct {
	coll *mongo.Collection
}

// NewTodoRepository returns a TodoRepository backed by the "todos" collection.
func NewTodoRepository(db *mongo.Database) *TodoRepository {
	return &TodoRepository{coll: db.Collection("todos")}
}

// Create inserts a new todo, populating its ID and timestamps.
func (r *TodoRepository) Create(ctx context.Context, todo *models.Todo) error {
	todo.ID = bson.NewObjectID()
	now := time.Now().UTC()
	todo.CreatedAt = now
	todo.UpdatedAt = now

	_, err := r.coll.InsertOne(ctx, todo)
	return err
}

// ListByUser returns all todos owned by userID, newest first.
func (r *TodoRepository) ListByUser(ctx context.Context, userID string) ([]models.Todo, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{"user_id": userID}, opts)
	if err != nil {
		return nil, err
	}

	todos := []models.Todo{}
	if err := cursor.All(ctx, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

// GetByID returns the todo with id, scoped to userID.
func (r *TodoRepository) GetByID(ctx context.Context, userID string, id bson.ObjectID) (*models.Todo, error) {
	var todo models.Todo
	err := r.coll.FindOne(ctx, bson.M{"_id": id, "user_id": userID}).Decode(&todo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// ToggleCompleted flips is_completed on the todo with id, scoped to userID,
// and returns the updated document.
func (r *TodoRepository) ToggleCompleted(ctx context.Context, userID string, id bson.ObjectID) (*models.Todo, error) {
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	pipeline := mongo.Pipeline{
		bson.D{{Key: "$set", Value: bson.D{
			{Key: "is_completed", Value: bson.D{{Key: "$not", Value: "$is_completed"}}},
			{Key: "updated_at", Value: time.Now().UTC()},
		}}},
	}

	var todo models.Todo
	err := r.coll.FindOneAndUpdate(ctx,
		bson.M{"_id": id, "user_id": userID},
		pipeline,
		opts,
	).Decode(&todo)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &todo, nil
}

// ResetCompleted sets is_completed to false for every todo owned by userID.
func (r *TodoRepository) ResetCompleted(ctx context.Context, userID string) error {
	_, err := r.coll.UpdateMany(ctx,
		bson.M{"user_id": userID},
		bson.M{"$set": bson.M{"is_completed": false, "updated_at": time.Now().UTC()}},
	)
	return err
}

// Delete removes the todo with id, scoped to userID.
func (r *TodoRepository) Delete(ctx context.Context, userID string, id bson.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id, "user_id": userID})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return util.ErrNotFound
	}
	return nil
}
