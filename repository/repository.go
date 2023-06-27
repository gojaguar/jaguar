package repository

import (
	"context"
)

// Repository holds the methods to interact with a persistence layer.
type Repository[T any] interface {
	Create(ctx context.Context, input T) (T, error)
	Find(ctx context.Context, query Query) ([]T, error)
	Get(ctx context.Context, query Query) (T, error)
	Update(ctx context.Context, query Query, data T) (T, error)
	Delete(ctx context.Context, query Query) (T, error)
}
