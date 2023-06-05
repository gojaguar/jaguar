package repository

import (
	"cloud.google.com/go/firestore"
	"context"
)

type firestoreRepository[T any] struct {
	collection string
	client     *firestore.Client
}

func (f *firestoreRepository[T]) Find(ctx context.Context, query Query) ([]T, error) {
	ref := f.client.Collection(f.collection)
	ref.Query = query.Firestore(ref.Query)

	snaps, err := ref.Documents(ctx).GetAll()
	if err != nil {
		return nil, err
	}

	var out []T
	var item T
	for _, s := range snaps {
		err = s.DataTo(item)
		if err != nil {
			return nil, err
		}
		out = append(out, item)
	}

	return out, nil
}

func NewRepositoryFirestore[T any](client *firestore.Client, collection string) Repository[T] {
	return &firestoreRepository[T]{
		client:     client,
		collection: collection,
	}
}
