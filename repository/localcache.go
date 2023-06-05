package repository

import (
	"context"
	"sync"
	"time"
)

type cached[T any] struct {
	CreatedAt time.Time
	Response  T
}

type localCacheRepository[T any] struct {
	repository Repository[T]
	cache      map[Query]cached[[]T]
	lock       sync.Mutex
	keepAlive  time.Duration
}

func (c *localCacheRepository[T]) Find(ctx context.Context, query Query) ([]T, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	item, ok := c.cache[query]
	now := time.Now()
	if ok {
		if item.CreatedAt.Add(c.keepAlive).After(now) {
			return item.Response, nil
		} else {
			delete(c.cache, query)
		}
	}
	response, err := c.repository.Find(ctx, query)
	if err != nil {
		return nil, err
	}
	c.cache[query] = cached[[]T]{
		CreatedAt: now,
		Response:  response,
	}
	return response, nil
}

func NewLocalCacheRepository[T any](r Repository[T], keepAlive time.Duration) Repository[T] {
	return &localCacheRepository[T]{
		repository: r,
		keepAlive:  keepAlive,
	}
}
