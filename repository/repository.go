package repository

import (
	"context"
)

// Repository holds the methods to interact with a persistence layer.
type Repository[T any] interface {
	// Find finds a list of entities using the criteria defined in the given Query.
	Find(ctx context.Context, query Query) ([]T, error)
}
