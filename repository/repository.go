package repository

import (
	"context"
)

// Repository contains a set of methods to interact with a persistence layer.
// The E parameter represents an entity model.
// The K parameter represent the primary key.
type Repository[E any, K comparable] interface {
	// Create creates an entity in a persistence layer.
	Create(ctx context.Context, entity E) (E, error)
	// CreateBulk creates a set of entities in a persistence layer.
	CreateBulk(ctx context.Context, entities []E) ([]E, error)
	// Get returns an entity from a persistence layer identified by its ID. It returns an error if the entity doesn't exist.
	Get(ctx context.Context, id K) (E, error)
	// Find returns a set of entities from a persistence layer identified by their ID. It returns
	// an empty slice if no records were found.
	Find(ctx context.Context, ids []K) ([]E, error)
	// Update updates an entity.
	Update(ctx context.Context, id K, entity E) (E, error)
	// UpdateBulk updates multiple entities with values of entity.
	UpdateBulk(ctx context.Context, ids []K, entity E) ([]E, error)
	// Remove removes the given id from the persistence layer.
	Remove(ctx context.Context, id K) (E, error)
	// RemoveBulk removes a set of elements from the persistence layer.
	RemoveBulk(ctx context.Context, ids []K) ([]E, error)
}
