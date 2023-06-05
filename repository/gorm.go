package repository

import (
	"context"
	"gorm.io/gorm"
)

type gormRepository[T any] struct {
	db *gorm.DB
}

func (g *gormRepository[T]) Find(ctx context.Context, query Query) ([]T, error) {
	q := query.GORM(g.db)
	var out []T
	if err := q.Model(g.model()).WithContext(ctx).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

func (g gormRepository[T]) model() any {
	var model T
	return model
}

func NewRepositoryGorm[T any](db *gorm.DB) Repository[T] {
	return &gormRepository[T]{
		db: db,
	}
}
