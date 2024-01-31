package repository

import (
	"context"
	"gorm.io/gorm"
)

// SQL implements Repository using gorm.
type SQL[E any, K comparable] struct {
	db *gorm.DB
}

// Create creates an entity in a persistence layer.
func (r *SQL[E, K]) Create(ctx context.Context, entity E) (E, error) {
	if err := r.db.WithContext(ctx).Model(new(E)).Create(&entity).Error; err != nil {
		var zero E
		return zero, err
	}
	return entity, nil
}

// CreateBulk creates a set of entities in a persistence layer.
func (r *SQL[E, K]) CreateBulk(ctx context.Context, entities []E) ([]E, error) {
	if err := r.db.WithContext(ctx).Model(new(E)).Create(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Get returns an entity from a persistence layer identified by its ID. It returns an error if the entity doesn't exist.
func (r *SQL[E, K]) Get(ctx context.Context, id K) (E, error) {
	var out E
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id = ?", id).First(&out).Error; err != nil {
		var zero E
		return zero, err
	}
	return out, nil
}

// Find returns a set of entities from a persistence layer identified by their IDs. It returns
// an empty slice if no records were found.
func (r *SQL[E, K]) Find(ctx context.Context, ids []K) ([]E, error) {
	var out []E
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id IN (?)", ids).Find(&out).Error; err != nil {
		return nil, err
	}
	return out, nil
}

// Update updates an entity.
func (r *SQL[E, K]) Update(ctx context.Context, id K, entity E) (E, error) {
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id = ?", id).Updates(&entity).Error; err != nil {
		var zero E
		return zero, err
	}
	return entity, nil
}

// UpdateBulk updates multiple entities with values of entity.
func (r *SQL[E, K]) UpdateBulk(ctx context.Context, ids []K, entity E) ([]E, error) {
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id IN (?)", ids).Updates(&entity).Error; err != nil {
		return nil, err
	}

	var result []E
	if err := r.db.WithContext(ctx).Model(new(E)).Where("id IN (?)", ids).Find(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// Remove removes the given id from the persistence layer.
func (r *SQL[E, K]) Remove(ctx context.Context, id K) (E, error) {
	entity, err := r.Get(ctx, id)
	if err != nil {
		var zero E
		return zero, err
	}

	if err := r.db.WithContext(ctx).Model(new(E)).Where("id = ?", id).Delete(&entity).Error; err != nil {
		var zero E
		return zero, err
	}

	return entity, nil
}

// RemoveBulk removes a set of elements from the persistence layer.
func (r *SQL[E, K]) RemoveBulk(ctx context.Context, ids []K) ([]E, error) {
	result, err := r.Find(ctx, ids)
	if err != nil {
		return nil, err
	}

	if err := r.db.WithContext(ctx).Model(new(E)).Where("id IN (?)", ids).Delete(&result).Error; err != nil {
		return nil, err
	}

	return result, nil
}

// NewRepositorySQL initializes a new implementation of Repository using an SQL ORM: gorm.
func NewRepositorySQL[E any, K comparable](db *gorm.DB) Repository[E, K] {
	return &SQL[E, K]{
		db: db,
	}
}
