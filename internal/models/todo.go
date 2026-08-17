// Package models contains shared data models / DTOs.
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// Todo is a single to-do item, owned by the Firebase user identified by UserID.
type Todo struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID      string        `bson:"user_id"        json:"-"`
	Title       string        `bson:"title"          json:"title"`
	IsCompleted bool          `bson:"is_completed"   json:"is_completed"`
	CreatedAt   time.Time     `bson:"created_at"     json:"created_at"`
	UpdatedAt   time.Time     `bson:"updated_at"     json:"updated_at"`
}

type TodoCreateRequest struct {
	Title string `json:"title"       binding:"required"`
}
